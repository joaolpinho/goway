package router
import (
	"fmt"
	"github.com/andrepinto/goway/product"
)

type RouterOptions func(*GoWayRouter) *GoWayRouter


//noinspection GoUnusedExportedFunction
func StrictSlash(r *GoWayRouter) *GoWayRouter {
	r.StrictSlash = true
	return r
}

//noinspection GoUnusedExportedFunction
func AddOptionsRoute(r *GoWayRouter) *GoWayRouter {
	r.AddOptionsRoute = true
	return r
}


type GoWayRouter struct {
	Router *Router
	AddOptionsRoute bool
	StrictSlash		bool
}

func NewGoWayRouter(opts ...RouterOptions) *GoWayRouter{
	r := &GoWayRouter{
		Router: NewRouter(),
	}
	for _, funcOpts := range opts {
		r = funcOpts(r)
	}
	return r
}


func (r *GoWayRouter) Compile()  {
	if (r.StrictSlash) {
		r.Router.CompileStrict()
	} else {
		r.Router.Compile()
	}
}


func (r *GoWayRouter) CheckRoute(path string, verb string, code string, version string) (*Route, map[string]interface{})  {
	return r.Router.Dispatch(verb, path, code, version)
}

func (r *GoWayRouter) CreateRoute(code string, version string, routes []product.Routes_v1)  {
	for _, k := range routes{
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
	}
}