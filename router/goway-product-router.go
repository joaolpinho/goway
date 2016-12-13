package router

import (
	"github.com/andrepinto/goway/product"
)

type GowayProductRouter struct  {
	GoWayRouter *GoWayRouter
	Products map[string]product.Product_v1
}

type GowayProductRouterOptions struct  {
	AddOptionsRoute bool
}

func NewGowayProductRouter(options *GowayProductRouterOptions) *GowayProductRouter{
	return &GowayProductRouter{
		GoWayRouter: NewGoWayRouter(options.AddOptionsRoute),
		Products: make(map[string]product.Product_v1),
	}
}


func (r *GowayProductRouter) LoadRoutes(products []product.Product_v1)  {
	for _, v := range products{
		r.Products[v.Code]=v
		r.GoWayRouter.CreateRoute(v.Code, v.Version, v.Routes)
	}


	r.GoWayRouter.Compile()
}


func (r *GowayProductRouter) CheckRoute(path string, verb string, code string, version string) (*Route, map[string]interface{})  {
	route, params := r.GoWayRouter.CheckRoute(path, verb, code, version)
	return route, params
}
