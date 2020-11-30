package sms_controller_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	url "net/url"
	"strings"
	"testing"
	"uims/app"
	"uims/boot"
	responses2 "uims/internal/controllers/responses"
	"uims/internal/controllers/sms_controller"
	requests2 "uims/internal/controllers/sms_controller/requests"
	"uims/pkg/tool"
)

var router *gin.Engine
var w *httptest.ResponseRecorder

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	router = app.GetEngineRouter()
	w = httptest.NewRecorder()
	m.Run()
}

func TestSendSMSverifyCode(t *testing.T) {
	request := requests2.SMSCodeByPhoneRequest{
		Phone: "13641337591",
		SPM:   "1024.YSyn5CeEqsVfEUqP.100.101",
	}
	requestStr, _ := json.Marshal(request)
	u := url.Values{}
	u.Add("phone", request.Phone)
	u.Add("spm", "1024.YSyn5CeEqsVfEUqP.100.101")
	req, _ := http.NewRequest("GET", "/api/sms/verifycode/send?"+u.Encode(), strings.NewReader(string(requestStr)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	response := responses2.Parse(w.Body.Bytes())
	assert.Equal(t, responses2.CodeSuccess, response.Code)
	tool.Dump(response)
}

func TestIsExistUserByScene(t *testing.T) {
	phone := "13641337591"
	exist, err := sms_controller.IsExistUserByScene(sms_controller.UsePhoneLogin, phone)
	assert.Nil(t, err)
	assert.Equal(t, true, exist)

	email := "342448932@qq.com"
	exist, err = sms_controller.IsExistUserByScene(sms_controller.UseEmailLogin, email)
	assert.Nil(t, err)
	assert.Equal(t, true, exist)
}

func TestCanSendSMS(t *testing.T) {
	phone := "13000000000"
	can, err := sms_controller.CanSendSMS(sms_controller.UsePhoneRegister, phone)
	assert.Nil(t, err)
	assert.Equal(t, true, can)

	can, err = sms_controller.CanSendSMS(sms_controller.UsePhoneLogin, phone)
	assert.EqualError(t, err, "用户未注册")
	assert.Equal(t, false, can)

	can, err = sms_controller.CanSendSMS(sms_controller.UsePhoneFindPasswd, phone)
	assert.EqualError(t, err, "用户未注册")
	assert.Equal(t, false, can)

	email := "3424489322@qq.com"
	can, err = sms_controller.CanSendSMS(sms_controller.UseEmailRegister, email)
	assert.Nil(t, err)
	assert.Equal(t, true, can)

	can, err = sms_controller.CanSendSMS(sms_controller.UseEmailLogin, email)
	assert.EqualError(t, err, "用户未注册")
	assert.Equal(t, false, can)

	can, err = sms_controller.CanSendSMS(sms_controller.UseEmailFindPasswd, email)
	assert.EqualError(t, err, "用户未注册")
	assert.Equal(t, false, can)
}
