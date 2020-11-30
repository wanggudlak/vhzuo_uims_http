package auth_controller_test

import (
	"encoding/json"
	"github.com/srlemon/gen-id"
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/app"
	"uims/boot"
	"uims/internal/controllers/auth_controller"
	"uims/internal/controllers/responses"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/encryption"
	"uims/pkg/test"
)

var httptest *test.Http

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	httptest = test.New(app.GetEngineRouter())
	m.Run()
}

func TestFindPasswordToken(t *testing.T) {
	user := &model.User{}
	err := db.Def().First(&user).Error
	assert.Nil(t, err)
	assert.NotEmpty(t, *user.Phone)
	resp := httptest.Get("/api/auth/find/password/token", auth_controller.FindPasswordTokenRequest{
		Phone:       *user.Phone,
		SMSCode:     "123456",
		SPMFullCode: "123456",
	})
	t.Logf("resp: %s", resp.Body)
	assert.Equal(t, resp.Code, 200)
	r := responses.Response{}
	err = json.Unmarshal(resp.Body.Bytes(), &r)
	if body, ok := r.Body.(map[string]interface{}); !ok {
		t.Error("响应处理失败", body)
		t.FailNow()
	} else {
		assert.NotEmpty(t, body["find_password_token"])
	}
}

func TestFindPassword(t *testing.T) {
	user := &model.User{}
	err := db.Def().First(&user).Error
	assert.Nil(t, err)
	assert.NotEmpty(t, *user.Phone)

	resp1 := httptest.Get("/api/auth/find/password/token", auth_controller.FindPasswordTokenRequest{
		Phone:       *user.Phone,
		SMSCode:     "123456",
		SPMFullCode: "123456",
	})
	assert.Equal(t, resp1.Code, 200)

	r1 := responses.Response{}
	err = json.Unmarshal(resp1.Body.Bytes(), &r1)
	assert.Nil(t, err)
	var findPasswordToken = ""
	if body, ok := r1.Body.(map[string]interface{}); !ok {
		t.Error("响应处理失败", body)
		t.FailNow()
	} else {
		assert.NotEmpty(t, body["find_password_token"])
		findPasswordToken = body["find_password_token"].(string)
	}

	resp := httptest.Post("/api/auth/find/password", auth_controller.FindPasswordRequest{
		FindPasswordToken: findPasswordToken,
		NewPassword:       "123456",
	})
	assert.Equal(t, resp.Code, 200)
}

func TestRegister(t *testing.T) {
	var err error
	var phone = genid.NewGeneratorData().PhoneNum
	w := httptest.Post("/api/auth/register", auth_controller.RegisterRequest{
		Phone:       phone,
		SMSCode:     "123456",
		Password:    "123456",
		SPMFullCode: "1024.ZBCDASDFGASDFASF.100.101",
	})
	t.Logf("resp: %s", w.Body)
	assert.Equal(t, w.Code, 200)
	r := responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	assert.Equal(t, 0, r.Code)

	// 检查用户与密码
	var user model.User
	err = db.Def().Where("phone = ?", phone).First(&user).Error
	assert.Nil(t, err)
	assert.True(t, encryption.BcryptCheck("123456", user.Passwd))
	var userInfo model.UserInfo
	err = db.Def().Where("phone = ?", phone).First(&userInfo).Error
	assert.Nil(t, err)
}
