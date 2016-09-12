package handlers
import (
	"github.com/andrepinto/goway/router"
	"net/http"
)


type HandlerFunc func(route *router.Route, r *http.Request)(bool)

type HandlerMap map[string]HandlerFunc


type IHandlerWorker interface {
	Run(handler string, route *router.Route) bool
	Add(action string, handler HandlerFunc)
}


type HandlerWorker struct{
	HandlerMap HandlerMap
}

func NewHandlerWorker() *HandlerWorker {
	return &HandlerWorker{HandlerMap{}}
}


func (hl *HandlerWorker) Add(action string, handler HandlerFunc) {
	hl.HandlerMap[action] = handler
}


func (hl *HandlerWorker) Run(handler string, route *router.Route, r *http.Request) bool{
	fn := hl.HandlerMap[handler]
	if fn == nil {
		return true
	}

	rs := fn(route, r)

	return rs
}