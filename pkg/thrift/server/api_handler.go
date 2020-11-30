package thriftserver

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
	"uims/gen-go/uims_rpc_api"
	"uims/pkg/glog"
	"uims/pkg/thrift/common"
	"uims/pkg/tool"
)

var Logger *logrus.Logger

type MethodFunc func(*Context)
type MethodMap map[string]MethodFunc

type BaseRequest struct {
	Method string      `json:"method_name"`
	Params interface{} `json:"params"`
}

func init() {
	Logger = glog.Channel("thrift")
}

func (req BaseRequest) JSON() string {
	return tool.JSONString(map[string]string{
		"method_name": req.Method,
		"params":      tool.JSONString(req.Params),
	})
}

type UIMSRpcAPIHandler struct {
	methodMap *MethodMap
}

func NewUIMSRpcAPIHandler() *UIMSRpcAPIHandler {
	return &UIMSRpcAPIHandler{}
}

func (handler *UIMSRpcAPIHandler) RegisterAPI(methodName string, fn MethodFunc) *UIMSRpcAPIHandler {
	if nil == handler.methodMap {
		handler.methodMap = &MethodMap{}
	}
	(*handler.methodMap)[methodName] = fn
	return handler
}

// 通过map来注册可调用方法
func (handler *UIMSRpcAPIHandler) RegisterAPIwhithMap(methodMap *MethodMap) *UIMSRpcAPIHandler {
	handler.methodMap = methodMap
	return handler
}

// params:
// {
//		"method_name": "example",
//		// json 字符串
//		"params": "{\"test\": \"test\"}"
// }
func (handler *UIMSRpcAPIHandler) InvokeMethod(_ context.Context, params string) (r *uims_rpc_api.Response, err error) {
	c := NewContext()
	logf("Call request: [%s]", params)
	func() {
		// 解析请求
		if err = c.ParseRequest(params); err != nil {
			logf("Invalid request: [%+v]", err)
			c.Response.SystemErr(err.Error())
			return
		}
		// 查找处理函数
		fn := (*handler.methodMap)[c.Request.MethodName]
		if nil == fn {
			logf("Not register method: [%s]", c.Request.MethodName)
			c.Response.SystemErr(common.INVALID_METHOD_NAME)
			return
		}
		// 执行
		fn(c)
		return
	}()
	resp := c.Response.ConvertThriftResp()
	logf("Call response: [%+v]", resp)
	return resp, nil
}

func logf(format string, v ...interface{}) {
	if v != nil {
		if Logger != nil {
			Logger.Infof(format, v...)
		} else {
			log.Printf(format, v...)
		}
	}
}
