package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"flag"
	"net/http"
	"os"
	"strings"
)

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8082", "proxy listen address")
	flag.Parse()

	os.Setenv("http_proxy", "http://localhost:8080")
	log.Println(getenvEitherCase("HTTP_PROXY"))
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		log.Println("Req")
		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		log.Println("Res")
		return resp
	})

	log.Fatal(http.ListenAndServe(*addr, proxy))

}

func getenvEitherCase(k string) string {
	if v := os.Getenv(strings.ToUpper(k)); v != "" {
		return v
	}
	return os.Getenv(strings.ToLower(k))
}