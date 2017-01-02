package router

import (
	"regexp"
	"strconv"
	"strings"
	"github.com/andrepinto/goway/product"
)


type Router struct {

	validExtensions []string
	routes map[string]*Route
	routesByMethod map[string][]*Route
	regexesByMethod map[string][]*regexp.Regexp

}


func NewRouter() *Router {
	router := new(Router)

	router.routes = make(map[string]*Route)

	// Allow requests with no extensions by default
	router.validExtensions = append(router.validExtensions, "")


	router.routesByMethod = map[string][]*Route{
		"HEAD":    make([]*Route, 0),
		"GET":    make([]*Route, 0),
		"POST":   make([]*Route, 0),
		"PUT":    make([]*Route, 0),
		"PATCH":    make([]*Route, 0),
		"DELETE": make([]*Route, 0),
		"OPTIONS": make([]*Route, 0),
	}

	router.regexesByMethod = map[string][]*regexp.Regexp{
		"HEAD":    make([]*regexp.Regexp, 0),
		"GET":    make([]*regexp.Regexp, 0),
		"POST":   make([]*regexp.Regexp, 0),
		"PUT":    make([]*regexp.Regexp, 0),
		"PATCH":    make([]*regexp.Regexp, 0),
		"DELETE": make([]*regexp.Regexp, 0),
		"OPTIONS": make([]*regexp.Regexp, 0),
	}

	return router
}

func (r *Router) Head(name string, pattern string, code string, version string, handlers []string, apiMethod product.Routes_v1 ) *Route {
	return r.addRoute("HEAD", name, pattern, code, version, handlers, apiMethod)
}

func (r *Router) Get(name string, pattern string, code string, version string, handlers []string, apiMethod product.Routes_v1 ) *Route {
	return r.addRoute("GET", name, pattern, code, version, handlers, apiMethod)
}

func (r *Router) Post(name string, pattern string, code string, version string, handlers []string, apiMethod product.Routes_v1) *Route {
	return r.addRoute("POST", name, pattern, code, version, handlers, apiMethod)
}

func (r *Router) Put(name string, pattern string, code string, version string, handlers []string, apiMethod product.Routes_v1) *Route {
	return r.addRoute("PUT", name, pattern, code, version, handlers, apiMethod)
}

func (r *Router) Patch(name string, pattern string, code string, version string, handlers []string, apiMethod product.Routes_v1) *Route {
	return r.addRoute("PATCH", name, pattern, code, version, handlers, apiMethod)
}

func (r *Router) Delete(name string, pattern string, code string, version string, handlers []string, apiMethod product.Routes_v1) *Route {
	return r.addRoute("DELETE", name, pattern, code, version, handlers, apiMethod)
}

func (r *Router) Options(name string, pattern string, code string, version string, handlers []string, apiMethod product.Routes_v1) *Route {
	return r.addRoute("OPTIONS", name, pattern, code, version, handlers, apiMethod)
}

func (r *Router) addRoute(method string, name string, pattern string, code string, version string,  handlers []string, apiMethod product.Routes_v1) *Route {


	route := newRoute(method, name, pattern, handlers, code, version, apiMethod)

	r.routes[name] = route
	r.routesByMethod[method] = append(r.routesByMethod[method], route)

	return route
}

func (r *Router) FindRoute(name string) (route *Route, found bool) {
	route, found = r.routes[name]

	return
}

func (r *Router) ValidExtensions(extensions ...string) *Router {
	r.validExtensions = extensions

	return r
}

func (r *Router) Compile() *Router {
	for method := range r.routesByMethod {
		pattern := ""
		for i, route := range r.routesByMethod[method] {
			pattern += "(?P<" + route.Name + ">/)" + strings.TrimRight(strings.TrimLeft(route.pattern, "/"),"/") + "(?:/)?|"
			if i > 0 && i%15 == 0 {
				r.regexesByMethod[method] = append(r.regexesByMethod[method], regexp.MustCompile("^(?:" + strings.TrimRight(pattern, "|") + ")$"))
				pattern = ""
				continue
			}
		}
		pattern = "^(?:" + strings.TrimRight(pattern, "|") + ")$"
		r.regexesByMethod[method] = append(r.regexesByMethod[method], regexp.MustCompile(pattern))
	}
	return r
}

func (r *Router) CompileStrict() *Router {
	for method := range r.routesByMethod {
		pattern := ""
		for i, route := range r.routesByMethod[method] {
			pattern += "(?P<" + route.Name + ">/)" + strings.TrimLeft(route.pattern, "/") + "|"
			if i > 0 && i%15 == 0 {
				r.regexesByMethod[method] = append(r.regexesByMethod[method], regexp.MustCompile("^(?:" + strings.TrimRight(pattern, "|") + ")$"))
				pattern = ""
				continue
			}
		}
		pattern = "^(?:" + strings.TrimRight(pattern, "|") + ")$"
		r.regexesByMethod[method] = append(r.regexesByMethod[method], regexp.MustCompile(pattern))
	}
	return r
}

func (r *Router) extensionIsValid(ext string) bool {
	for _, valid := range r.validExtensions {
		if ext == valid {
			return true
		}
	}

	return false
}

func (r *Router) Dispatch(method string, path string, code string, version string) (*Route, map[string]interface{}) {

	params := make(map[string]interface{})
	var ext string
	var match []string
	var rrt *Route



	if !r.extensionIsValid(ext) {
		return nil, nil
	}



	for _, compiled := range r.regexesByMethod[method] {
		if match = compiled.FindStringSubmatch(path); match == nil {
			continue
		}

		for i, name := range compiled.SubexpNames() {

			nm := strings.Split(name, "_")

			if(len(nm)<2){
				continue
			}

			if(nm[0]!=version || nm[1]!=code){
				continue
			}

			paramLength := len(params)
			if i == 0 || match[i] == "" {
				if paramLength == 0 {
					continue
				}

				// All Params have been set. Empty matches means all params have been captured
				break
			}



			if paramLength == 0 {
				// Capture the name and set the ext so len(params) returns 1 on next loop
				rrt ,_ = r.FindRoute(name)

				params["ext"] = ext

				continue
			}

			if intValue, err := strconv.Atoi(match[i]); err == nil {
				params[name] = intValue

				continue
			}

			params[name] = match[i]
		}

		return rrt, params
	}

	return nil, nil
}


