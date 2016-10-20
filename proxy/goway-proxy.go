package proxy

import (
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/url"
	"net/http"
	"net/http/httputil"

	. "github.com/andrepinto/goway"
	"github.com/andrepinto/goway/router"
	"github.com/andrepinto/goway/product"
	"github.com/andrepinto/goway/util"
	"github.com/andrepinto/goway/handlers"
	"github.com/andrepinto/goway/util/worker"
)


type GoWayProxy struct{
	proxy        	 	*httputil.ReverseProxy
	target        		*url.URL
	productRouter       	*router.GowayProductRouter
	clientRouter        	*router.GowayClientRouter
	handlerWorker		*handlers.HandlerWorker
	HttpRequestLog 		HttpRequestLog
	TaskWorker 	        worker.ITaskWorker
}

type GowayProxyOptions struct {
	Target 			string
	ProductRouter 		*router.GowayProductRouter
	ClientRouter 		*router.GowayClientRouter
	HandlerWorker 		*handlers.HandlerWorker
	TaskWorker 		worker.ITaskWorker
}

//noinspection GoUnusedExportedFunction
func NewGoWayProxy(options *GowayProxyOptions) *GoWayProxy{
	target, _ := url.Parse(options.Target)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &transport{http.DefaultTransport}

	return &GoWayProxy{
		proxy: proxy,
		//target: target,
		productRouter	: options.ProductRouter,
		clientRouter	: options.ClientRouter,
		handlerWorker	: options.HandlerWorker,
		TaskWorker	: options.TaskWorker,
	}
}


func (p *GoWayProxy) Handle(w http.ResponseWriter, req *http.Request) {

	var rs bool
	var route *router.Route

	res :=  NewHttpResponse(w)
	version := req.Header.Get(GOWAY_VERSION)
	if ( version == "" ) {
		req.Header.Set(GOWAY_VERSION, DEFAULT_VERSION)
		version = DEFAULT_VERSION
	}


	rs, cl, newPath := p.checkClient(req.URL.Path, version)
	req.URL.Path = newPath


	if(!rs) {
		p.respond(req, res.Set( http.StatusNotFound, API_KEY_NOT_FOUND, nil) )
		return
	}

	//check client routes
	rs, route = p.checkRoute(newPath, req.Method, cl.Client, cl.Version, true)
	if(rs){
		p.redirect(route, cl.GlobalInjectData, req, res, cl.Product, cl.Client, cl.Version)
		return
	}

	//check product routes
	rs, route = p.checkRoute(newPath, req.Method, cl.Product, cl.Version, false)
	if(rs){
		p.redirect(route, cl.GlobalInjectData, req, res, cl.Product, cl.Client, cl.Version)
		return
	}

	p.respond(req, res.Set(http.StatusNotFound, API_ROUTE_NOT_FOUND, nil) )
}

func(p *GoWayProxy) checkRoute(path string, verb string, code string, version string, client bool) (bool, *router.Route){
	var route *router.Route;

	if(client){
		route, _ = p.clientRouter.CheckRoute(path, verb, code, version)
	}else{
		route, _ = p.productRouter.CheckRoute(path, verb, code, version)
	}


	if(route==nil){
		return false, nil
	}else{
		return true, route
	}
}

func(p *GoWayProxy) checkClient(path string, version string) (bool, *product.Client_v1, string){
	urlSplit := strings.Split(path, "/")

	if(len(urlSplit)==0){
		return false, nil, ""
	}

	client := p.clientRouter.CheckClient(urlSplit[1], version)

	if(client==nil || len(client.Client)==0){
		return false, client, ""
	}

	urlWithoutApiId := fmt.Sprintf("/%s",strings.Join(urlSplit[2:],"/"))

	return true, client, urlWithoutApiId
}

func(p *GoWayProxy) respond( req *http.Request, res *HttpResponse ) {

	response := res.Dispatch( req.Header.Get("Accept") )
	end := time.Now()


	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	request := buf.String()

	log := LogRecord{

		Time:          	end.UTC(),
		Ip:            	strings.Split(req.RemoteAddr, ":")[0],
		Method:        	req.Method,
		Uri:           	req.RequestURI,
		Username:      	"",
		Protocol:      	req.Proto,
		Host:          	req.Host,
		Status:        	res.Status,
		Size:          	int64(len(response)),
		ElapsedTime:   	end.Sub(res.StartTime),
		RequestHeader: 	req.Header,
		ResBody:		response,
		ReqBody: 		request,
		ServicePath:   	req.URL.Path,
		Product:       	req.Header.Get(GOWAY_PRODUCT),
		Client:        	req.Header.Get(GOWAY_CLIENT),
		Version:       	req.Header.Get(GOWAY_VERSION),

	}

	opt := map[string]string{}
	job := worker.Job{Name: REQUEST_LOGGER_EMMIT, Resource: nil, Payload:log, Map:opt, Id:""}
	worker.JobQueue <- job
}

func(p *GoWayProxy) redirect(route *router.Route, globalInjectData []product.InjectData_v1, req *http.Request, res *HttpResponse, product string, client string, version string) {

	if(route.ApiMethod.InjectGlobalData){
		p.injectDataValues(util.MergeInjectData(globalInjectData,route.ApiMethod.InjectData), req)
	}else{
		p.injectDataValues(route.ApiMethod.InjectData, req)
	}


	err := p.dispatchHandlers(route, req)
	if(err != nil){
		p.respond( req, res.Set( err.Status, err.Message, err.Data ) )
		return
	}

	req.URL.Path = fmt.Sprintf("%s%s", route.ApiMethod.ServiceName, req.URL.Path)

	req.Header.Add(GOWAY_PRODUCT, product)
	req.Header.Add(GOWAY_CLIENT, client)
	req.Header.Add(GOWAY_VERSION, version)

	res.ResponseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	p.proxy.ServeHTTP(res.ResponseWriter, req)

}

func(p *GoWayProxy) injectDataValues(values []product.InjectData_v1, r *http.Request){
	for _, v := range values{
		if(v.Where==product.WHERE_HEADER){
			//w.Header().Set(v.Code, v.Value)
			r.Header.Add(v.Code, v.Value)
			continue
		}

		if(v.Where==product.WHERE_PARAMS){
			values := r.URL.Query()
			values.Add(v.Code, v.Value)
			r.URL.RawQuery = values.Encode()
			continue
		}

		if(v.Where==product.WHERE_URL){
			r.URL.Path = fmt.Sprintf("/%s/%s%s", v.Code, v.Value, r.URL.Path)
			continue
		}

	}
}

func(p *GoWayProxy) dispatchHandlers(route *router.Route, req *http.Request) (*handlers.HandlerError){

	for _, v := range route.ApiMethod.Handlers{
		response := p.handlerWorker.Run(v, route, req)
		if(response != nil){
			return response
		}
	}

	return nil
}