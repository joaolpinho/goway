package logger

import (
	"net/http"
)


type LogDestination interface {
	Debug(data string, args ...interface{})
	Info(data string, args ...interface{})
	Warning(data string, args ...interface{})
	Error(data string, args ...interface{})
	Critical(data string, args ...interface{})
}

type Formatter interface {
	Format(interface{}) (string, error)
}

type HTTPLogger struct {
	name        string
	formatter   Formatter
	destination LogDestination
}

func (log *HTTPLogger) format(res http.ResponseWriter, req *http.Request) (string, error) {
	return log.formatter.Format(newFields(res, req))
}

// NewHTTPLogger constructor
func NewHTTPLogger(name string, destination LogDestination, formatter Formatter) (*HTTPLogger, error) {
	return &HTTPLogger{name: name, formatter: formatter, destination: destination}, nil
}

func (log *HTTPLogger) Debug(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Debug(data)
}

func (log *HTTPLogger) Info(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Info(data)
}

func (log *HTTPLogger) Warning(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Warning(data)
}

func (log *HTTPLogger) Error(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Error(data)
}

func (log *HTTPLogger) Critical(res http.ResponseWriter, req *http.Request) {
	data, _ := log.format(res, req)
	log.destination.Critical(data)
}
