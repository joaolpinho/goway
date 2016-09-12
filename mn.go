package main

import (
	"net/http"
	"flag"
	"github.com/andrepinto/goway/proxy"
	"fmt"
	"github.com/andrepinto/goway/router"
	"github.com/andrepinto/goway/product"
	"github.com/andrepinto/goway/handlers"
	jwt "github.com/andrepinto/goway/handlers/jwt"

	"github.com/andrepinto/goway/custom"
	"github.com/andrepinto/goway/util/worker"
)

var httpRequestLog proxy.HttpRequestLog

func main()  {

	fmt.Println("Started")

	var (
		port 		= flag.String("port", ":8084", "8084")
		maxWorkers   	= flag.Int("max_workers", 20, "The number of workers to start")
		maxQueueSize 	= flag.Int("max_queue_size", 20, "The size of job queue")
		url 		= flag.String("url", "http://localhost:9000", "http://localhost:9000")
	)

	flag.Parse()

	/* -------- INIT WORKERS --------- */

	worker.JobQueue = make(chan worker.Job, *maxQueueSize)

	dispatcher := worker.NewDispatcher(worker.JobQueue, *maxWorkers)

	taskWork := worker.NewTaskWorker()

	dispatcher.Run(taskWork)

	/* -------- END WORKERS --------- */


	productResource := product.NewProductResource(&product.ProductResourceOptions{})

	gowayProductRouter := router.NewGowayProductRouter()

	gowayProductRouter.LoadRoutes(productResource.GetAllProducts())

	gowayClientRouter := router.NewGowayClientRouter()

	gowayClientRouter.LoadRoutes(productResource.GetAllClients())

	handlersWork := handlers.NewHandlerWorker()
	handlersWork.Add("AUTHENTICATION", jwt.Jwt)

	//httpRequestLog := proxy.NewBasicLog()
	httpRequestLog = custom.NewElasticLog("http://52.30.6.179:9200", "gateway","http-logger")
	httpRequestLog.Start()

	gowayProxy := proxy.NewGoWayProxy(&proxy.GowayProxyOptions{
		Target		: 	*url,
		ProductRouter	: 	gowayProductRouter,
		ClientRouter	:   	gowayClientRouter,
		HandlerWorker	:  	handlersWork,
		TaskWorker	: 	taskWork,
	})


	gowayProxy.TaskWorker.AddJob(proxy.REQUEST_LOGGER_EMMIT, SendHttpLogs)


	// server
	http.HandleFunc("/", gowayProxy.Handle)

	http.ListenAndServe(*port, nil)


}


func SendHttpLogs(job *worker.Job)(bool) {

	fmt.Println(job)

	log := job.Payload.(proxy.LogRecord)

	httpRequestLog.Log(&log)

	return true;
}