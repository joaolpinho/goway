package router

import (
	"regexp"
	"strings"
	"github.com/andrepinto/goway/product"
)


type Route struct {
	method string
	Name string
	pattern string
	basePattern string
	Handlers []string
	Code string
	Version string
	ApiMethod product.Routes_v1
}


func newRoute(method string, name string, pattern string, handlers []string, code string, version string, apiMethod product.Routes_v1) (route *Route) {
	pattern = strings.Replace(pattern, "(int)", "([0-9]+)", -1)
	pattern = strings.Replace(pattern, "(alpha)", "([a-z]+)", -1)
	pattern = strings.Replace(pattern, "(alphanumeric)", "([a-z0-9]+)", -1)
	pattern = strings.Replace(pattern, "(slug)", "([a-z0-9-]+)", -1)
	pattern = strings.Replace(pattern, "(mongo)", "([0-9a-fA-F]{24})", -1)
	pattern = strings.Replace(pattern, "(md5)", "([0-9a-fA-F]{32})", -1)

	named := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	namedRegex := regexp.MustCompile(`:([a-zA-Z0-9_]+)\(([^\)]+)\)`)

	// Create basic base pattern used for URL generation
	basePattern := namedRegex.ReplaceAllString(pattern, ":$1")

	if match := namedRegex.MatchString(pattern); match {
		pattern = namedRegex.ReplaceAllString(pattern, "(?P<$1>$2)")
	} else {
		pattern = named.ReplaceAllString(pattern, "(?P<$1>[^/]+)")
	}

	route = &Route{strings.ToUpper(method), name, pattern, basePattern, handlers, code, version, apiMethod}

	return
}
