package logger

import "fmt"

type BasicLog struct {

}

func NewBasicLog() *BasicLog{
	return &BasicLog{}
}


func (log *BasicLog) Debug(data string, args ...interface{}){
 	fmt.Println(data, args)
}
func (log *BasicLog) Info(data string, args ...interface{}){
	fmt.Println(data, args)
}
func (log *BasicLog) Warning(data string, args ...interface{}){
	fmt.Println(data, args)
}
func (log *BasicLog) Error(data string, args ...interface{}){
	fmt.Println(data, args)
}
func (log *BasicLog) Critical(data string, args ...interface{}){
	fmt.Println(data, args)
}
