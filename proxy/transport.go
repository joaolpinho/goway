package proxy

import (
	"io/ioutil"
	"net/http"
	"time"
	"strings"
	"bytes"
	"strconv"
	"github.com/andrepinto/goway/util/worker"
)

type transport struct {
	http.RoundTripper
}


func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	startTime := time.Now()

	ip := strings.Split(req.RemoteAddr, ":")[0]



	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {

	}


	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	log := LogRecord{

			Time:          	startTime.UTC(),
			Ip:            	ip,
			Method:        	req.Method,
			Uri:           	req.RequestURI,
			Username:      	"",
			Protocol:      	req.Proto,
			Host:          	req.Host,
			Status:        	resp.StatusCode,
			Size:          	int64(len(resp.Body)),
			ElapsedTime:   	time.Duration(0),
			RequestHeader: 	req.Header,
			ResBody:		string(b[:]),
			ServicePath:   	resp.Request.URL.Path,
			Product:       	req.Header.Get(GOWAY_PRODUCT),
			Client:        	req.Header.Get(GOWAY_CLIENT),
			Version:       	req.Header.Get(GOWAY_VERSION),

	}

	finishTime := time.Now()

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	request := buf.String()
	log.ReqBody = request

	log.Time = finishTime.UTC()
	log.ElapsedTime = finishTime.Sub(startTime)


	opt := map[string]string{}


	job := worker.Job{Name: REQUEST_LOGGER_EMMIT, Resource: nil, Payload:log, Map:opt, Id:""}
	worker.JobQueue <- job


	return resp, nil
}
