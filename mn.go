package main

import (
	"net/http"
	"flag"
	"github.com/andrepinto/goway/proxy"
	"fmt"
	"github.com/andrepinto/goway/router"
	"github.com/andrepinto/goway/product"
)

func main()  {

	fmt.Println("Started")

	port := flag.String("port", ":8081", "8081")
	url := flag.String("url", "http://localhost:8080", "http://localhost:8080")

	flag.Parse()

	productResource := product.NewProductResource(&product.ProductResourceOptions{})

	gowayRouter := router.NewGoWayRouter()

	gowayRouter.LoadProductRoutes(productResource.GetAllProducts())

	gowayRouter.LoadClients(productResource.GetAllClients())


	gowayProxy := proxy.NewGoWayProxy(*url, gowayRouter)

	// server
	http.HandleFunc("/", gowayProxy.Handle)
	http.ListenAndServe(*port, nil)




}
