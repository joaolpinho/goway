package router
import (
	"fmt"
	"github.com/andrepinto/goway/product"
	"github.com/andrepinto/goway/util"
)

type GoWayRouter struct {
	Router *Router
	Clients map[string]product.Client_v1
	Products map[string]product.Product_v1
}

func NewGoWayRouter() *GoWayRouter{
	return &GoWayRouter{
		Router: NewRouter(),
		Clients: make(map[string]product.Client_v1),
		Products: make(map[string]product.Product_v1),
	}
}

/*
ROUTES
 */

func (r *GoWayRouter) LoadProductRoutes(products []product.Product_v1){

	for _, v := range products{
		r.Products[v.Code]=v
		r.AddProductRoutes(v)
	}


	r.Compile()
}

func (r *GoWayRouter) Compile(){
	r.Router.Compile()
}

func (r *GoWayRouter) CheckRoute(path string, verb string, product string, version string, client string) (*Route, map[string]interface{}) {
	route, params := r.Router.Dispatch(verb, path, product, version, client)
	return route, params

}

func (r *GoWayRouter) AddProductRoutes(product product.Product_v1){
	for _, k := range product.Routes{
		r.AddRoute(fmt.Sprintf("%s_%s_%s",product.Version, product.Code, k.Code), k.ListenPath, k.Verb,product.Code, product.Version, k.Handlers, k)
	}
}

func (r *GoWayRouter) AddRoute(name string, path string, verb string,  product string, version string, handlers []string, apiMethod product.Routes_v1){

	switch verb {
	case "GET":
		r.Router.Get(name, path, product, version, handlers, apiMethod)
	case "POST":
		r.Router.Post(name, path, product, version, handlers, apiMethod)
	case "PUT":
		r.Router.Put(name, path, product, version, handlers, apiMethod)
	case "DELETE":
		r.Router.Delete(name, path, product, version, handlers, apiMethod)
	default:


	}

}

/*
CLIENTS
 */

func (r *GoWayRouter) LoadClients(clients []product.Client_v1){
	for _, k := range clients{
		r.Clients[util.ClientCode(k.ApiPath, k.Version)] = k
	}
}

func (r *GoWayRouter) CheckClient(path string, version string) *product.Client_v1{
	x:= r.Clients[util.ClientCode(path, version)]
	return &x
}

