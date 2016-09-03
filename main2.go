package main

import (
	"github.com/andrepinto/goway/router"
	"fmt"
)

func main(){
	router := router.NewRoute()
	router.Get("read", "/articles/:name")
	router.Compile()

	route, params := router.Dispatch("GET", "/articles/article-title")
	route2, params2 := router.Dispatch("GET", "/articles/10")
	route3, params3 := router.Dispatch("GET", "/photos/10")

	fmt.Println(route.Name, route.URL(), params)
	fmt.Println(route2.Name, route2.URL(), params2)
	fmt.Println(route3, route3, params3)
}
