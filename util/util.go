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
