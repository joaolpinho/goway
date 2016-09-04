package proxy

import(
	"encoding/json"
)

type HttpResponse struct {
	Status	int	`json:"status"`
	Message string	`json:"message"`
}

func NewHttpResponse(status int, message string) string{

	msg := HttpResponse{
		Status: status,
		Message: message,
	}

	b, err := json.Marshal(msg)

	if err != nil {
		return ""
	}

	return string(b)

}