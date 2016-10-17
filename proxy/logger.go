package proxy

import (
	"net/http"
	"time"
)


type LogRecord struct {
	Time                                      time.Time 		`json:"time"`
	Ip					  string		`json:"ip"`
	Method 					  string		`json:"method"`
	Uri					  string 		`json:"uri"`
	Protocol 				  string		`json:"protocol"`
	Username 				  string		`json:"username"`
	Host 					  string		`json:"host"`
	ServicePath				  string		`json:"service_path"`
	Product 				  string		`json:"product"`
	Client 					  string		`json:"client"`
	Version 				  string		`json:"version"`
	Status                                    int			`json:"status"`
	Size                                      int64			`json:"size"`
	ElapsedTime                               time.Duration		`json:"elapsed_time"`
	RequestHeader                             http.Header		`json:"request_header"`
	CustomRecords                             map[string]string	`json:"custom_records"`
	ReqBody					  string		`json:"request_body"`
	ResBody					  string		`json:"response_body"`

}

type HttpRequestLog interface {
	Log(record interface{})
	Start() error
}