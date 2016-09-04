package handlers
import "github.com/andrepinto/goway/router"


type HandlerFunc func(route *router.Route)(bool)

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


func (hl *HandlerWorker) Run(handler string, route *router.Route) bool{
	fn := hl.HandlerMap[handler]
	if fn == nil {
		return false
	}

	fn(route)

	return true
}