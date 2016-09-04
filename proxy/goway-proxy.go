package proxy

import (
	"fmt"
	"github.com/andrepinto/goway/router"
	"github.com/andrepinto/goway/product"
	"strings"
	"github.com/andrepinto/goway/util"
	"github.com/andrepinto/goway/handlers"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	DEFAULT_VERSION = "1"
	GOWAY_PRODUCT= "GOWAY_PRODUCT"
	GOWAY_CLIENT= "GOWAY_CLIENT"
	GOWAY_VERSION = "GOWAY_VERSION"
	API_NOT_FOUND = "API_NOT_FOUND"
	ROUTE_NOT_FOUND = "ROUTE_NOT_FOUND"
	AUTHENTICATION_HANDLER = "AUTHENTICATION"
	CUSTOM_HANDLER = "HANDLER_ERROR"
)

type GoWayProxy struct{
	proxy        	 	*httputil.ReverseProxy
	target        		*url.URL
	productRouter       	*router.GowayProductRouter
	clientRouter        	*router.GowayClientRouter
	handlerWorker		*handlers.HandlerWorker
	HttpRequestLog 		HttpRequestLog
}

func NewGoWayProxy(target string, productRouter *router.GowayProductRouter, clientRouter *router.GowayClientRouter, handlerWorker *handlers.HandlerWorker, httpRequestLog HttpRequestLog) *GoWayProxy{
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &transport{http.DefaultTransport, httpRequestLog}

	return &GoWayProxy{
		proxy: proxy,
		//target: target,
		productRouter: productRouter,
		clientRouter: clientRouter,
		handlerWorker: handlerWorker,
		HttpRequestLog: httpRequestLog,
	}
}


func (p *GoWayProxy) Handle(w http.ResponseWriter, r *http.Request) {

	//change
	version := DEFAULT_VERSION
	var rs bool
	var route *router.Route

	rs, cl, newPath := p.checkClient(r.URL.Path, version)

	r.URL.Path = newPath


	if(!rs){
		http.Error(w, NewHttpResponse(http.StatusNotFound, API_NOT_FOUND), http.StatusNotFound)
		return
	}

	//check client routes
	rs, route = p.checkRoute(newPath, r.Method, cl.Client, cl.Version, true)

	if(rs){
		p.redirect(route, cl.GlobalInjectData, w, r, cl.Product, cl.Client, cl.Version)
	}else{
		//check product routes
		rs, route = p.checkRoute(newPath, r.Method, cl.Product, cl.Version, false)

		if(rs){
			p.redirect(route, cl.GlobalInjectData, w, r, cl.Product, cl.Client, cl.Version)
		}else{

			http.Error(w, NewHttpResponse(http.StatusNotFound, API_NOT_FOUND), http.StatusNotFound)
			return
		}
	}



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

func(p *GoWayProxy) redirect(route *router.Route, globalInjectData []product.InjectData_v1, w http.ResponseWriter, r *http.Request, product string, client string, version string){

	if(route.ApiMethod.InjectGlobalData){
		p.injectDataValues(util.MergeInjectData(globalInjectData,route.ApiMethod.InjectData), r)
	}else{
		p.injectDataValues(route.ApiMethod.InjectData, r)
	}


	result, v :=p.dispatchHandlers(route, r)
	if(!result){
		if(v==AUTHENTICATION_HANDLER){
			http.Error(w, NewHttpResponse(http.StatusUnauthorized, AUTHENTICATION_HANDLER), http.StatusUnauthorized)
			return
		}else{
			http.Error(w, NewHttpResponse(http.StatusBadRequest, CUSTOM_HANDLER), http.StatusBadRequest)
			return
		}
	}

	r.URL.Path = fmt.Sprintf("%s%s", route.ApiMethod.ServiceName, r.URL.Path)

	r.Header.Add(GOWAY_PRODUCT, product)
	r.Header.Add(GOWAY_CLIENT, client)
	r.Header.Add(GOWAY_VERSION, version)

	p.proxy.ServeHTTP(w, r)

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

func(p *GoWayProxy) dispatchHandlers(route *router.Route, r *http.Request)(bool, string){
	for _, v := range route.ApiMethod.Handlers{
		result := p.handlerWorker.Run(v, route, r)
		if(!result){
			return result, v
		}
	}

	return true, ""

}
