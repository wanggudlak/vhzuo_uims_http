package thriftserver

import (
	"encoding/json"
	"fmt"
	//json "github.com/json-iterator/go"
	"strings"
	validator2 "uims/internal/validator"
	"uims/pkg/thrift/common"
)

// thrift 请求响应上下文处理
// 支持功能: validate 结构体验证
// 简化响应

type handlerFunc func(ctx *Context)

type Context struct {
	Request  *Request
	Response *Response
	// 用以支持中间件
	handlers []handlerFunc
	// 处理器处理到了第几个方法
	index int8
}

func NewContext() *Context {
	return &Context{
		Response: NewResponse(),
		Request:  NewRequest(),
	}
}

func (c *Context) ParseRequest(body string) error {
	var err error
	type params struct {
		Method string `json:"method_name"`
		Params string `json:"params"`
	}
	var p params
	err = json.Unmarshal([]byte(body), &p)
	if err != nil {
		// 解析失败(非可以解析的 JSON 字符串)
		return err
	}
	c.Request.Body = p.Params
	c.Request.MethodName = p.Method
	if c.Request.MethodName == "" {
		return fmt.Errorf("%s", common.INVALID_METHOD_NAME)
	}
	return nil
}

func (c *Context) ShouldBind(obj interface{}) error {
	if c.Request == nil || c.Request.Body == "" {
		return fmt.Errorf("invalid request")
	}
	decoder := json.NewDecoder(strings.NewReader(c.Request.Body))
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}

func validate(obj interface{}) error {
	if validator2.Validate == nil {
		return nil
	}
	return validator2.Validate.Struct(obj)
}
