package thriftserver

import (
	"fmt"
	"github.com/mailru/easyjson/buffer"
	"gopkg.in/go-playground/validator.v9"
	"uims/gen-go/uims_rpc_api"
	validator2 "uims/internal/validator"
	"uims/pkg/glog"
	"uims/pkg/thrift/common"
	"uims/pkg/tool"
)

type Biz struct {
	Content interface{} `json:"biz_content"`
	Status  string      `json:"biz_status"`
	Message string      `json:"biz_message"`
}

type Response struct {
	Message string `json:"message"`
	Data    Biz
	Status  string `json:"status"`
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) SystemErr(msg string) {
	r.Message = msg
	r.Status = common.STATUS_FAILED
}

// 结构体绑定失败时将解析错误消息
// 使用该方法会将参数验证的第一个错误消息作为 biz_message 返回
// 使用该方法会自动翻译错误信息
func (r *Response) BadParams(err error) {
	if err == nil {
		return
	}
	r.Message = common.CALL_SUCCESS_MSG
	r.Status = common.STATUS_SUCCESS

	if invalid, ok := err.(*validator.InvalidValidationError); ok {
		r.Data = Biz{
			Content: "",
			Status:  common.STATUS_PARAMS_FAILED,
			Message: fmt.Sprintf("参数解析错误: %s", invalid.Error()),
		}
		return
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errMsg buffer.Buffer
		validationErrorsTranslations := validationErrors.Translate(validator2.Trans)
		// 取第一个错误信息
		for _, value := range validationErrorsTranslations {
			errMsg.AppendString(value)
			break
		}
		r.Data = Biz{
			Content: "",
			Status:  common.STATUS_PARAMS_FAILED,
			Message: string(errMsg.Buf),
		}
		return
	}
	glog.Channel("thrift").WithError(err).Error("错误的绑定信息")
	r.SystemErr(common.INVALID_PARAMS)
	return
}

// 处理业务失败
func (r *Response) Fail(msg string) {
	r.Message = common.CALL_SUCCESS_MSG
	r.Status = common.STATUS_SUCCESS
	r.Data = Biz{
		Content: "",
		Status:  common.STATUS_FAILED,
		Message: msg,
	}
}

// 处理业务成功
func (r *Response) Success(data interface{}, msg string) {
	r.Message = common.CALL_SUCCESS_MSG
	r.Status = common.STATUS_SUCCESS
	r.Data = Biz{
		Content: data,
		Status:  common.STATUS_SUCCESS,
		Message: msg,
	}
}

// 处理错误请求
// 根据需求, 这里可以判断 error 的类型
func (r *Response) Error(err error) {
	r.Fail(err.Error())
}

// 响应结构应为 (看接口示例是这样写的, 就这样转换)
// {
// 		// 系统级错误
// 		"msg": "",
//		"status": "",
//		"data": {
// 			// 业务级数据
//			"biz_content": "",
//			"biz_status": "",
// 			"biz_message: ""
//		},
// }
func (r *Response) ConvertThriftResp() *uims_rpc_api.Response {
	return &uims_rpc_api.Response{
		Status: r.Status,
		Msg:    r.Message,
		Data:   tool.JSONString(r.Data),
	}
}
