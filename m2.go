package main

import (
	"fmt"
	"net/http"
	logger "github.com/andrepinto/goway/util/logger"
)

var logs = logger.NewBasicLog()

func main() {

	// Boot web server and listen on 8080
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8888", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	formatter, _ := logger.NewTextFormatter("myHandler")
	clerk, err := logger.NewHTTPLogger("myHandler", logs, formatter)
	if err != nil {
		fmt.Println("HTTP logger could not be created", err)
	}
	defer clerk.Info(w, r)

	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])

}
