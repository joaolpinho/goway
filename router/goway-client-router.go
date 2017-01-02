package router

import (
	"github.com/andrepinto/goway/product"
	"github.com/andrepinto/goway/util"
	"github.com/andrepinto/goway/constants"
)

type GowayClientRouter struct  {
	Clients map[string]*InternalClientRouter
}

type InternalClientRouter struct {
	Client *product.Client_v1
	Router   *GoWayRouter
}

//noinspection GoUnusedExportedFunction
func NewGowayClientRouter() *GowayClientRouter{
	r := &GowayClientRouter{
		map[string]*InternalClientRouter{},
	}
	return r
}


func (r *GowayClientRouter) AddClient(client product.Client_v1, mode string, filters []string, options ...RouterOptions){
	internalRouter := &InternalClientRouter{
		&client,
		NewGoWayRouter(options...),
	}

	internalRouter.Router.CreateRoute(util.ClientRouteCode(client.Client, client.Product), client.Version, util.FilterClientRoutesByAssets(&client, filters,  util.FilterByAsset))
	internalRouter.Router.Compile()

	if mode==constants.CLIENT_HEADERS_MODE {
		r.Clients[util.ClientApiHeaders(client.Client, client.Product, client.Version)]=internalRouter
	}else {
		r.Clients[util.ClientApiKey(client.ApiPath, client.Version)]=internalRouter
	}



}

func (r *GowayClientRouter) CheckClientByApiKey(key string, version string) *InternalClientRouter{
	x:= r.Clients[util.ClientApiKey(key, version)]
	return x
}

func (r *GowayClientRouter) CheckClientByHeaders(client string, product string, version string) *InternalClientRouter{
	x:= r.Clients[util.ClientApiHeaders(client, product, version)]
	return x
}