package login_controller_test

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/app"
	"uims/boot"
	"uims/conf"
	"uims/internal/controllers/login_controller/contexts"
	"uims/internal/controllers/responses"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/test"
)

var httptest *test.Http

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	httptest = test.New(app.GetEngineRouter())
	conf.Switch.CSRF = false
	m.Run()
}

func TestAuthenticateCassAdminLogin(t *testing.T) {
	reqBody := contexts.AccountPasswdVerifyCodeAuthRequest{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		Account:         "admin",
		Password:        "123456",
		Code:            "123456",
		VerificationKey: "1234567890",
	}

	w := httptest.Post("/api/login/authenticate", reqBody)
	t.Logf("resonse: %s", w.Body)
	assert.Equal(t, 200, w.Code)
	r := responses.Response{}
	err := json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.NotEmpty(t, body["redirect_url"])
		assert.Contains(t, body["redirect_url"], "code")
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}
}

func TestAuthenticateCassWebLogin(t *testing.T) {
	// 测试第一步登录
	reqBody := contexts.AccountPasswdVerifyCodeSMSCodeAuthRequest{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    1,
			AuthScene:   2,
			SPMFullCode: "1024.YSyn5CeEqsVfEUqP.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		Account:         "admin",
		Password:        "123456",
		Code:            "123456",
		VerificationKey: "1234567890",
		PasswordToken:   "",
		SMSCode:         "",
	}

	// 第一步登录
	w := httptest.Post("/api/login/authenticate", reqBody)
	t.Logf("resonse: %s", w.Body)
	assert.Equal(t, 200, w.Code)
	type body struct {
		PasswordToken string `json:"password_token"`
	}
	type resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Body    body   `json:"body"`
	}

	r := resp{}
	err := json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	assert.Equal(t, 0, r.Code)
	assert.Equal(t, "需要进行手机验证", r.Message)
	assert.NotEmpty(t, r.Body.PasswordToken)

	reqBody.PasswordToken = r.Body.PasswordToken

	// 未填写短信验证码
	w = httptest.Post("/api/login/authenticate", reqBody)
	t.Logf("response2: %s", w.Body)
	type rs2 struct {
		responses.Response
	}
	r2 := rs2{}
	err = json.Unmarshal(w.Body.Bytes(), &r2)
	assert.Nil(t, err)
	assert.Equal(t, 1, r2.Code)
	assert.Equal(t, "短信验证码必须传入", r2.Message)

	// 填写短信验证码
	reqBody.SMSCode = "123456"
	w = httptest.Post("/api/login/authenticate", reqBody)
	t.Logf("response3: %s", w.Body)
	r3 := responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r3)
	assert.Nil(t, err)
	assert.Equal(t, 0, r3.Code)
	if body, ok := r3.Body.(map[string]interface{}); ok {
		assert.NotEmpty(t, body["redirect_url"])
		assert.Contains(t, body["redirect_url"], "code")
	} else {
		t.Error("响应格式不正确", r3.Body)
		t.FailNow()
	}
}

func TestEmptyQuery(t *testing.T) {
	userWeChat := model.UserWeChat{}
	err := db.Def().
		Where(&model.UserWeChat{
			WeChatOpenId:  "test",
			WeChatUnionId: "",
		}).
		First(&userWeChat).Error
	t.Logf("userwechat: %+v", userWeChat)
	assert.True(t, gorm.IsRecordNotFoundError(err))
}

func TestFreezeUser(t *testing.T) {
	err := db.Def().Model(&model.User{}).Where("account = ?", "admin").UpdateColumn("status", "N").Error
	assert.Nil(t, err)

	reqBody := contexts.AccountPasswdVerifyCodeAuthRequest{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		Account:         "admin",
		Password:        "123456",
		Code:            "123456",
		VerificationKey: "1234567890",
	}

	w := httptest.Post("/api/login/authenticate", reqBody)
	err = db.Def().Model(&model.User{}).Where("account = ?", "admin").UpdateColumn("status", "Y").Error
	assert.Nil(t, err)
	t.Logf("resonse: %s", w.Body)
	assert.Equal(t, 200, w.Code)
	r := responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	assert.Equal(t, 1, r.Code)
	assert.Equal(t, "用户已冻结, 无法登录", r.Message)
}
