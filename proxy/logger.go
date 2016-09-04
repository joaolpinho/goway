package proxy

import (
	"net/http"
	"time"
)


type LogRecord struct {
	Time                                      time.Time
	Ip, Method, Uri, Protocol, Username, Host string
	ServicePath				  string
	Product 				  string
	Client 					  string
	Version 				  string
	Status                                    int
	Size                                      int64
	ElapsedTime                               time.Duration
	RequestHeader                             http.Header
	CustomRecords                             map[string]string
	Body					  string

}

