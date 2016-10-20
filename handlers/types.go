package handlers

type HandlerError struct {
	XMLName 			struct{} 				`json:"-" xml:"root"`
	Status				int						`json:"status" xml:"status"`
	Message 			string					`json:"message" xml:"message"`
	Data				interface{}				`json:"result,omitempty" xml:"result,omitempty"`
}
func NewHttpError( status int, message string ) *HandlerError {
	return &HandlerError{
		Status: status,
		Message: message,
	}
}
func ( e *HandlerError ) Error( ) string {
	return e.Message
}