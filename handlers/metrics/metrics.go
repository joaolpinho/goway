package handlers
import (
	"github.com/andrepinto/goway/router"
	"fmt"
	"log"
	"net/http"
	"time"
)


func Metrics(route *router.Route)(bool){
	fmt.Println("metrics")
	return true
}


func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}