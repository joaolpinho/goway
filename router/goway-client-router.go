package router

import (
	"github.com/andrepinto/goway/product"
	"github.com/andrepinto/goway/util"
)

type GowayClientRouterOptions struct  {
	AddOptionsRoute bool
}

type GowayClientRouter struct  {
	GoWayRouter *GoWayRouter
	Clients map[string]product.Client_v1
}


func NewGowayClientRouter(options *GowayClientRouterOptions) *GowayClientRouter{
	return &GowayClientRouter{
		GoWayRouter: NewGoWayRouter(options.AddOptionsRoute),
		Clients: make(map[string]product.Client_v1),
	}
}


func (r *GowayClientRouter) LoadRoutes(clients []product.Client_v1)  {
	for _, v := range clients{
		r.Clients[util.ClientCode(v.ApiPath, v.Version)]=v
		r.GoWayRouter.CreateRoute(v.Client, v.Version, v.Routes)
	}


	r.GoWayRouter.Compile()
}


func (r *GowayClientRouter) CheckRoute(path string, verb string, code string, version string) (*Route, map[string]interface{})  {
	route, params := r.GoWayRouter.CheckRoute(path, verb, code, version)
	return route, params
}


func (r *GowayClientRouter) CheckClient(path string, version string) *product.Client_v1{
	x:= r.Clients[util.ClientCode(path, version)]
	return &x
}