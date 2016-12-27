package router

import (
	"github.com/andrepinto/goway/product"
	"github.com/andrepinto/goway/util"
	"github.com/andrepinto/goway/proxy"
)

type GowayClientRouter struct  {
	*GoWayRouter
	Clients map[string]product.Client_v1
}


//noinspection GoUnusedExportedFunction
func NewGowayClientRouter(options ...RouterOptions) *GowayClientRouter{
	r := &GowayClientRouter{
		NewGoWayRouter(options...),
		map[string]product.Client_v1{},
	}
	return r
}


func (r *GowayClientRouter) LoadRoutes(clients []product.Client_v1, mode string)  {
	for _, v := range clients{
		if mode==proxy.CLIENT_HEADERS_MODE {
			r.Clients[util.ClientApiHeaders(v.Client, v.Product, v.Version)]=v
		}else {
			r.Clients[util.ClientApiKey(v.ApiPath, v.Version)]=v
		}
		r.GoWayRouter.CreateRoute(v.Client, v.Version, v.Routes)
	}


	r.GoWayRouter.Compile()
}

func (r *GowayClientRouter) CheckClientByApiKey(path string, version string) *product.Client_v1{
	x:= r.Clients[util.ClientApiKey(path, version)]
	return &x
}

func (r *GowayClientRouter) CheckClientByHeaders(client string, product string, version string) *product.Client_v1{
	x:= r.Clients[util.ClientApiHeaders(client, product, version)]
	return &x
}