package util

import (
	"fmt"
	"github.com/andrepinto/goway/product"
)

func ClientCode(path string, version string) string{
	return fmt.Sprint("[%s-%s]", version, path)
}

func MergeInjectData(global []product.InjectData_v1, method []product.InjectData_v1) []product.InjectData_v1{
	result := method

	if(len(global)==0){
		return method
	}

	for _, v := range global{
		for _, k := range method{
			if(v.Code==k.Code){
				break
			}
		}

		result = append(result, v)
	}

	return result
}


func FilterClientRoutesByAsset(cl []product.Client_v1, asset string, f func(product.Routes_v1, string) bool) []product.Client_v1 {
	clients := make([]product.Client_v1, 0)
	for _, c := range cl {
		routes := make([]product.Routes_v1, 0)
		for _, v := range c.Routes {
			if f(v, asset) {
				routes = append(routes, v)
			}
		}

		c.Routes = routes
		clients = append(clients, c)
	}


	return clients
}

func FilterByAsset(route product.Routes_v1, asset string) bool{
	return route.Asset==asset
}