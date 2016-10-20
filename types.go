package goway

import(
	"strings"
	"time"
	"net/http"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type HttpResponse struct {
	XMLName 			struct{} 				`json:"-" xml:"root"`
	StartTime			time.Time				`json:"-" xml:"-"`
	ResponseWriter  	http.ResponseWriter		`json:"-" xml:"-"`

	Status				int						`json:"status" xml:"status"`
	Success 			bool					`json:"success" xml:"success"`
	Message 			string					`json:"message" xml:"message"`
	Data				interface{}				`json:"result,omitempty" xml:"result,omitempty"`
}
func NewHttpResponse( w http.ResponseWriter ) *HttpResponse{
	return &HttpResponse{
		ResponseWriter: w,
		StartTime: time.Now(),
		Status: http.StatusOK,
		Message: "Ok",
	}
}
func ( res *HttpResponse ) Set( status int, message string, data interface{} ) *HttpResponse {
	res.Status = status
	res.Message = message
	res.Data = data

	return res
}
func ( res *HttpResponse ) Dispatch( mime string ) string {

	var body []byte

	w := res.ResponseWriter
	mime = strings.SplitN( mime, ";", 1 )[0]

	res.Success = res.Status >= 200 && res.Status < 300

	switch mime {
	case "application/json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		body, _ = json.Marshal(res)
		break

	case "text/xml":
	case "application/xml":
		w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", mime))
		body, _ = xml.Marshal(res)
		break
	default:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		body = []byte(res.Message)
		break
	}

	w.WriteHeader(res.Status)
	w.Write(body)
	return string(body)
}