package main

import (
	"net/http"
	"flag"
	"github.com/andrepinto/goway/proxy"
	"fmt"
	"github.com/andrepinto/goway/router"
	"github.com/andrepinto/goway/product"
	"github.com/andrepinto/goway/handlers"
	metrics "github.com/andrepinto/goway/handlers/metrics"

)

func main()  {

	fmt.Println("Started")

	port := flag.String("port", ":8084", "8084")
	url := flag.String("url", "http://localhost:9000", "http://localhost:9000")

	flag.Parse()

	productResource := product.NewProductResource(&product.ProductResourceOptions{})

	gowayProductRouter := router.NewGowayProductRouter()

	gowayProductRouter.LoadRoutes(productResource.GetAllProducts())

	gowayClientRouter := router.NewGowayClientRouter()

	gowayClientRouter.LoadRoutes(productResource.GetAllClients())

	handlersWork := handlers.NewHandlerWorker()
	handlersWork.Add("METRICS", metrics.Metrics)

	gowayProxy := proxy.NewGoWayProxy(*url, gowayProductRouter, gowayClientRouter, handlersWork)

	// server
	http.HandleFunc("/", gowayProxy.Handle)

	http.ListenAndServe(*port, nil)




}

