package proxy

import (
	"log"
	"encoding/json"
)

type BasicLog struct {

}

func NewBasicLog() *BasicLog{
	return &BasicLog{}
}

func(lg *BasicLog) Log(record *LogRecord){
	b, err := json.Marshal(record)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println(string(b))
}
