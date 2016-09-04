package proxy

import (
	"io/ioutil"
	"net/http"
	"time"
	"strings"
	"bytes"
	"strconv"
)

type transport struct {
	http.RoundTripper
	httpRequestLog HttpRequestLog
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
		return nil, err
	}

	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	log := &LogRecord{

			Time:          startTime.UTC(),
			Ip:            ip,
			Method:        req.Method,
			Uri:           req.RequestURI,
			Username:      "",
			Protocol:      req.Proto,
			Host:          req.Host,
			Status:        resp.StatusCode,
			Size:          0,
			ElapsedTime:   time.Duration(0),
			RequestHeader: req.Header,
			Body:	       string(b[:]),
			ServicePath:   resp.Request.URL.Path,
			Product:       req.Header.Get(GOWAY_PRODUCT),
			Client:        req.Header.Get(GOWAY_CLIENT),
			Version:       req.Header.Get(GOWAY_VERSION),

	}

	finishTime := time.Now()

	log.Time = finishTime.UTC()
	log.ElapsedTime = finishTime.Sub(startTime)

	t.httpRequestLog.Log(log)



	return resp, nil
}
