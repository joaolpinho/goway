package router

import (
	"github.com/andrepinto/goway/product"
	"github.com/andrepinto/goway/util"
	"fmt"
)

type GowayProductRouter struct  {
	Products map[string]*InternalProductRouter
}

type InternalProductRouter struct {
	Product *product.Product_v1
	Router   *GoWayRouter
}

//noinspection GoUnusedExportedFunction
func NewGowayProductRouter() *GowayProductRouter{
	r := &GowayProductRouter{
		map[string]*InternalProductRouter{},
	}
	return r
}

func (r *GowayProductRouter) AddProduct(product product.Product_v1, mode string, filters []string, options ...RouterOptions){
	internalRouter := &InternalProductRouter{
		&product,
		NewGoWayRouter(options...),
	}

	internalRouter.Router.CreateRoute(product.Code, product.Version, util.FilterProductsRoutesByAssets(&product, filters,  util.FilterByAsset))
	internalRouter.Router.Compile()

	r.Products[util.ProductCode(product.Code, product.Version)]=internalRouter


	fmt.Println("------------------------------")


}


func (r *GowayProductRouter) CheckProduct(code string, version string) *InternalProductRouter{
	x:= r.Products[util.ProductCode(code, version)]
	return x
}