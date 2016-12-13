package router
import (

	"github.com/andrepinto/goway/product"
	"fmt"
	"github.com/andrepinto/goway/util"
)

type GoWayRouter struct {
	Router *Router
	AddOptionsRoute bool

}

func NewGoWayRouter(addOptionsRoute bool) *GoWayRouter{
	return &GoWayRouter{
		Router: NewRouter(),
		AddOptionsRoute: addOptionsRoute,
	}
}


func (r *GoWayRouter) Compile()  {
	r.Router.Compile()
}


func (r *GoWayRouter) CheckRoute(path string, verb string, code string, version string) (*Route, map[string]interface{})  {
	path = util.PrepareUrl(path)
	route, params := r.Router.Dispatch(verb, path, code, version)
	return route, params
}

func (r *GoWayRouter) CreateRoute(code string, version string, routes []product.Routes_v1)  {
	for _, k := range routes{
		k.ListenPath = util.PrepareUrl(k.ListenPath)
		r.AddRoute(fmt.Sprintf("%s_%s_%s", version, code, k.Code), k.ListenPath, k.Verb, code, version, k.Handlers, k)
		if(len(k.Alias)>0){
			r.AddRoute(fmt.Sprintf("%s_%s_%s", version, code, fmt.Sprintf("%s_alias",k.Code)), k.Alias, k.Verb, code, version, k.Handlers, k)
		}
		if(r.AddOptionsRoute){
			r.AddRoute(fmt.Sprintf("%s_%s_%s", version, code, k.Code), k.ListenPath, "OPTIONS", code, version, k.Handlers, k)
		}
	}
}


func (r *GoWayRouter) AddRoute(name string, path string, verb string,  code string, version string, handlers []string, apiMethod product.Routes_v1){

	switch verb {
	case "GET":
		r.Router.Get(name, path, code, version, handlers, apiMethod)
	case "POST":
		r.Router.Post(name, path, code, version, handlers, apiMethod)
	case "PUT":
		r.Router.Put(name, path, code, version, handlers, apiMethod)
	case "DELETE":
		r.Router.Delete(name, path, code, version, handlers, apiMethod)
	case "OPTIONS":
		r.Router.Options(name, path, code, version, handlers, apiMethod)
	default:


	}

}



/*
CLIENTS


func (r *GoWayRouter) LoadClients(clients []product.Client_v1){
	for _, k := range clients{
		r.Clients[util.ClientCode(k.ApiPath, k.Version)] = k
	}
}

func (r *GoWayRouter) CheckClient(path string, version string) *product.Client_v1{
	x:= r.Clients[util.ClientCode(path, version)]
	return &x
}
*/
