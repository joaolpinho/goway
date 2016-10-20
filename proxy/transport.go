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


func (t *transport) RoundTrip(req *http.Request) (res *http.Response, err error) {



	reqBodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body.Close()

	reqBody := ioutil.NopCloser(bytes.NewReader(reqBodyBytes))
	req.Body = reqBody



	startTime := time.Now()

	res, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	finishTime := time.Now()



	ip := strings.Split(req.RemoteAddr, ":")[0]

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	resBody := ioutil.NopCloser(bytes.NewReader(resBodyBytes))
	res.Body = resBody
	res.ContentLength = int64(len(resBodyBytes))
	res.Header.Set("Content-Length", strconv.Itoa(len(resBodyBytes)))


	log := LogRecord{

		Time:           finishTime.UTC(),
		Ip:            	ip,
		Method:        	req.Method,
		Uri:           	req.RequestURI,
		Username:      	"",
		Protocol:      	req.Proto,
		Host:          	req.Host,
		Status:        	res.StatusCode,
		Size:          	res.ContentLength,
		ElapsedTime:   	finishTime.Sub(startTime),
		RequestHeader: 	req.Header,
		ReqBody:		string(reqBodyBytes),
		ResBody:		string(resBodyBytes),
		ServicePath:   	res.Request.URL.Path,
		Product:       	req.Header.Get(GOWAY_PRODUCT),
		Client:        	req.Header.Get(GOWAY_CLIENT),
		Version:       	req.Header.Get(GOWAY_VERSION),

	}

	opt := map[string]string{}
	job := worker.Job{Name: REQUEST_LOGGER_EMMIT, Resource: nil, Payload:log, Map:opt, Id:""}
	worker.JobQueue <- job


	return res, nil
}
