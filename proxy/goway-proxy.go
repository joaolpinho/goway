package proxy

import (
	"net/http/httputil"
	"net/http"
	"net/url"
	"fmt"
	"github.com/andrepinto/goway/router"
	"github.com/andrepinto/goway/product"
	"strings"
	"github.com/andrepinto/goway/util"
)

const (
	DEFAULT_VERSION = "1"
)

type GoWayProxy struct{
	proxy         *httputil.ReverseProxy
	target        *url.URL
	router        *router.GoWayRouter
}

func NewGoWayProxy(target string, router *router.GoWayRouter) *GoWayProxy{
	url, _ := url.Parse(target)


	return &GoWayProxy{
		proxy: httputil.NewSingleHostReverseProxy(url),
		target: url,
		router: router,
	}
}


func (p *GoWayProxy) Handle(w http.ResponseWriter, r *http.Request) {

	//change
	version := DEFAULT_VERSION

	rs, cl, newPath := p.checkClient(r.URL.Path, version)

	fmt.Println(rs, cl, newPath)

	if(!rs){
		http.Error(w, "API_NOT_FOUND", http.StatusNotFound)
		return
	}

	//check custom client routes
	if(len(cl.Routes)>0){

	}

	rs, route := p.checkRoute(newPath, r.Method, cl.Product, cl.Version, cl.Client)

	if(rs){
		p.redirect(route, cl.GlobalInjectData, w, r)
	}else{
		http.Error(w, "ROUTE_NOT_FOUND", http.StatusNotFound)
		return
	}

}

func(p *GoWayProxy) checkRoute(path string, verb string, product string, version string, client string) (bool, *router.Route){

	route, _ := p.router.CheckRoute(path, verb, product, version, client)
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

	client := p.router.CheckClient(urlSplit[1], version)

	if(client==nil || len(client.Client)==0){
		return false, client, ""
	}

	urlWithoutApiId := fmt.Sprintf("/%s",strings.Join(urlSplit[2:],"/"))

	return true, client, urlWithoutApiId
}

func(p *GoWayProxy) redirect(route *router.Route, globalInjectData []product.InjectData_v1, w http.ResponseWriter, r *http.Request){

	if(route.ApiMethod.InjectGlobalData){
		p.injectDataValues(util.MergeInjectData(globalInjectData,route.ApiMethod.InjectData), r)
	}else{
		p.injectDataValues(route.ApiMethod.InjectData, r)
	}


	p.dispatchHandlers(route)

	r.URL.Path = fmt.Sprintf("%s%s", route.ApiMethod.ServiceName, r.URL.Path)

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

func(p *GoWayProxy) dispatchHandlers(oute *router.Route){
	fmt.Println("RUN DUMMY HANDLER")
}