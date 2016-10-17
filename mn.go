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
	"github.com/andrepinto/goway/util/worker"
	"github.com/andrepinto/goway-elasticsearch-logging"
	_ "encoding/json"
	"github.com/andrepinto/goway-mongodb-store"
)

var httpRequestLog proxy.HttpRequestLog

func main()  {

	fmt.Println("Started")

	var (
		port 		= flag.String("port", ":8084", "8084")
		maxWorkers   	= flag.Int("max_workers", 20, "The number of workers to start")
		maxQueueSize 	= flag.Int("max_queue_size", 20, "The size of job queue")
		url 		= flag.String("url", "http://localhost:9000", "http://localhost:9000")
		elasticUrl 	= flag.String("elasticUrl", "http://52.30.6.179:9200", "http://52.30.6.179:9200")
		elasticIndex 	= flag.String("elasticIndex", "gateway", "gateway")
		elasticType 	= flag.String("elasticType", "http-logger", "http-logger")
	)

	flag.Parse()

	/* -------- INIT WORKERS --------- */

	worker.JobQueue = make(chan worker.Job, *maxQueueSize)

	dispatcher := worker.NewDispatcher(worker.JobQueue, *maxWorkers)

	taskWork := worker.NewTaskWorker()

	dispatcher.Run(taskWork)


	/* -------- INIT REPOSITORY --------- */

	//_ = goway_couchbase_store.NewLocalRepository()
	//
	//repoCouch := goway_couchbase_store.NewCouchbaseRepository(&goway_couchbase_store.CouchbaseRepositoryOptions{
	//	ConnectionString: "couchbase://52.30.6.179",
	//	BucketName: "gateway",
	//})

	//repoCouch.CreateAndGet();

	repoMongo := goway_mongodb_store.NewMongodbRepository(&goway_mongodb_store.MongodbRepositoryOptions{
		Url:"localhost:27017",
		DatabaseName:"goway",
	});

	//repoMongo.Create()
	fmt.Println(repoMongo.GetAllProducts())

	/* -------- INIT GOWAY --------- */


	productResource := product.NewProductResource(&product.ProductResourceOptions{
		Repository: repoMongo,
	})

	gowayProductRouter := router.NewGowayProductRouter()

	//b, _ := json.Marshal(productResource.GetAllProducts())
	//fmt.Println(b)


	gowayProductRouter.LoadRoutes(productResource.GetAllProducts())

	gowayClientRouter := router.NewGowayClientRouter()

	gowayClientRouter.LoadRoutes(productResource.GetAllClients())

	handlersWork := handlers.NewHandlerWorker()
	handlersWork.Add(proxy.AUTHENTICATION_HANDLER, jwt.Jwt)

	//httpRequestLog := proxy.NewBasicLog()
	httpRequestLog = goway_elasticsearch_logging.NewElasticLog(*elasticUrl, *elasticIndex, *elasticType)
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