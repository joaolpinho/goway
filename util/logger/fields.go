package logger

import (
	"net/http"
	"strconv"
)

type fields struct {
	Method  string              `json:"method"`
	Status  string              `json:"status"`
	Path    string              `json:"path"`
	Host    string              `json:"host"`
	Headers map[string][]string `json:"headers"`
}

func newFields(res http.ResponseWriter, req *http.Request) *fields {
	// If you need to add to the @fields key, add it here
	return &fields{
		Method:  req.Method,
		Status:  fetchStatusCode(res), // Only fetches if Status() is defined on res
		Path:    req.URL.Path,
		Headers: map[string][]string(req.Header),
		Host:    req.Host,
	}
}


func fetchStatusCode(res http.ResponseWriter) string {
	var statusCode int

	type statusInterface interface {
		Status() int
	}

	statusCaller, ok := res.(statusInterface)
	if ok {
		statusCode = statusCaller.Status()
	}

	if statusCode == 0 {
		return ""
	}

	return strconv.Itoa(statusCode)
}