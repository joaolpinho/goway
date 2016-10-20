package handlers
import (
	"net/http"
	"github.com/andrepinto/goway/router"
)


type HandlerFunc func(route *router.Route, r *http.Request)(*HandlerError)

type HandlerMap map[string]HandlerFunc


type IHandlerWorker interface {
	Run(handler string, route *router.Route) *HandlerError
	Add(action string, handler HandlerFunc)
}


type HandlerWorker struct{
	HandlerMap HandlerMap
}

//noinspection GoUnusedExportedFunction
func NewHandlerWorker() *HandlerWorker {
	return &HandlerWorker{HandlerMap{}}
}


func (hl *HandlerWorker) Add(action string, handler HandlerFunc) {
	hl.HandlerMap[action] = handler
}


func (hl *HandlerWorker) Run(handler string, route *router.Route, r *http.Request) *HandlerError {
	fn := hl.HandlerMap[handler]
	if fn == nil {
		return nil
	}

	return fn(route, r)
}