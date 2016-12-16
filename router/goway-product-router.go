package router

import (
	"github.com/andrepinto/goway/product"
)

type GowayProductRouter struct  {
	*GoWayRouter
	Products map[string]product.Product_v1
}

//noinspection GoUnusedExportedFunction
func NewGowayProductRouter( options ...RouterOptions) *GowayProductRouter{
	r := &GowayProductRouter{
		NewGoWayRouter(options...),
		map[string]product.Product_v1{},
	}
	return r
}


func (r *GowayProductRouter) LoadRoutes(products []product.Product_v1)  {
	for _, v := range products{
		r.Products[v.Code]=v
		r.GoWayRouter.CreateRoute(v.Code, v.Version, v.Routes)
	}


	r.GoWayRouter.Compile()
}
