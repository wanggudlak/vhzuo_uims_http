package thriftclient

import "encoding/json"

type TRequest interface {
	String() string
}

type BRequest struct {
	MethodName string      `json:"method_name"`
	Params     interface{} `json:"params"`
}

type Request struct {
	MethodName string      `json:"method_name"`
	Params     interface{} `json:"params"`
}

func (b BRequest) String() string {
	p2, _ := json.Marshal(b)
	return string(p2)
}
