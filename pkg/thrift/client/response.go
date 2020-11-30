package thriftclient

import (
	"encoding/json"
	"fmt"
	"uims/gen-go/uims_rpc_api"
)

type TResponse interface {
	OK() bool
	CallOK() bool
	Parse([]byte) error
	SetStatus(s string)
	SetMsg(s string)
	GetStatus() string
	ParseContent(t interface{}) error
}

type Biz struct {
	BizStatus  string      `json:"biz_status"`
	BizMessage string      `json:"biz_message"`
	BizContent interface{} `json:"biz_content"`
}

type BResponse struct {
	uims_rpc_api.Response
	Biz Biz `json:"biz"`
}

func (b *BResponse) Parse(bytes []byte) error {
	fmt.Println(string(bytes))
	err := json.Unmarshal(bytes, &b)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(b.Data), &b.Biz)
}

func (b *BResponse) ParseContent(t interface{}) error {
	bytes, err := json.Marshal(b.Biz.BizContent)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &t)
}

func (b *BResponse) SetStatus(s string) {
	b.Status = s
}
func (b *BResponse) SetMsg(s string) {
	b.Msg = s
}

func (b *BResponse) GetStatus() string {
	return b.Status
}

func (b *BResponse) CallOK() bool {
	return b.Status == "success"
}

func (b *BResponse) OK() bool {
	return b.Status == "success" && b.Biz.BizStatus == "success"
}

func (b *BResponse) Err() string {
	if b.Status != "success" {
		return b.Msg
	}
	if b.Biz.BizStatus != "success" {
		return b.Biz.BizMessage
	}
	return ""
}
