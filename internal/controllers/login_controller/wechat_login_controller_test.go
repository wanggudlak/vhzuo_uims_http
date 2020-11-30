package login_controller_test

import (
	"encoding/json"
	genid "github.com/srlemon/gen-id"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"uims/internal/controllers/login_controller/contexts"
	"uims/internal/controllers/responses"
	"uims/internal/service/wechatqrcode"
	"uims/pkg/randc"
)

func TestWeChatQRCode(t *testing.T) {
	reqBody := contexts.WeChatQRCodeReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
	}

	w := httptest.Get("/api/login/wechat/qr_code", reqBody)
	t.Logf("resonse: %s", w.Body)
	assert.Equal(t, 200, w.Code)
	r := responses.Response{}
	err := json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.NotEmpty(t, body["qr_code"])
		assert.Contains(t, body["qr_code"], "mp.weixin.qq.com")
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}
}

func TestWeChatQRCodeLogin(t *testing.T) {
	reqBody := contexts.WeChatQRCodeLoginReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		SceneId: 1,
	}

	w := httptest.Get("/api/login/wechat/qr_code/login", reqBody)
	t.Logf("resonse: %s", w.Body)
	assert.Equal(t, 200, w.Code)
	r := responses.Response{}
	err := json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.True(t, body["timeout"].(bool))
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}

	// 获取了二维码, 但是没有扫码
	sceneID := rand.Intn(999999999)
	err = wechatqrcode.SetScene(sceneID)

	reqBody = contexts.WeChatQRCodeLoginReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		SceneId: sceneID,
	}
	w = httptest.Get("/api/login/wechat/qr_code/login", reqBody)
	t.Logf("resonse2: %s", w.Body)
	r = responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.False(t, body["timeout"].(bool))
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}
}

func TestWeChatQRCodeLoginWithUserExists(t *testing.T) {
	var err error
	// 获取了二维码, 但是没有扫码
	sceneID := rand.Intn(999999999)
	err = wechatqrcode.SetScene(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetScanOK(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetSceneAuthOK(sceneID, 1)
	assert.Nil(t, err)
	reqBody := contexts.WeChatQRCodeLoginReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		SceneId: sceneID,
	}
	w := httptest.Get("/api/login/wechat/qr_code/login", reqBody)
	t.Logf("resonse2: %s", w.Body)
	r := responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.True(t, body["auth_ok"].(bool))
		assert.NotEmpty(t, body["redirect_url"].(string))
		assert.Contains(t, body["redirect_url"].(string), "http")
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}
}

func TestWeChatQRCodeLoginWithNeedBindPhone(t *testing.T) {
	var err error
	// 获取了二维码, 但是没有扫码
	sceneID := rand.Intn(999999999)
	err = wechatqrcode.SetScene(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetScanOK(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetNeedBindPhone(sceneID, wechatqrcode.WeChatInfo{
		City:       "武汉市",
		Province:   "湖北省",
		Sex:        1,
		Country:    "中国",
		HeadImgURL: "",
		Nickname:   "小白",
		OpenID:     randc.RandStringN(28),
	})
	assert.Nil(t, err)
	reqBody := contexts.WeChatQRCodeLoginReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		SceneId: sceneID,
	}
	w := httptest.Get("/api/login/wechat/qr_code/login", reqBody)
	t.Logf("resonse2: %s", w.Body)
	r := responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.False(t, body["auth_ok"].(bool))
		assert.True(t, body["need_bind_phone"].(bool))
		assert.NotEmpty(t, body["bind_phone_token"].(string))
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}
}

func TestWeChatQRCodeLoginBindPhone(t *testing.T) {
	var err error
	sceneID := rand.Intn(999999999)
	err = wechatqrcode.SetScene(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetScanOK(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetNeedBindPhone(sceneID, wechatqrcode.WeChatInfo{
		City:       "武汉市",
		Province:   "湖北省",
		Sex:        1,
		Country:    "中国",
		HeadImgURL: "",
		Nickname:   "小白",
		OpenID:     randc.RandStringN(28),
	})
	assert.Nil(t, err)
	reqBody := contexts.WeChatQRCodeLoginReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		SceneId: sceneID,
	}
	w := httptest.Get("/api/login/wechat/qr_code/login", reqBody)
	t.Logf("resonse2: %s", w.Body)
	r := responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	bindPhoneToken := ""
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.False(t, body["auth_ok"].(bool))
		assert.True(t, body["need_bind_phone"].(bool))
		assert.NotEmpty(t, body["bind_phone_token"].(string))
		bindPhoneToken = body["bind_phone_token"].(string)
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}

	// 绑定
	rb := contexts.BindPhoneReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		Phone:          genid.NewGeneratorData().PhoneNum,
		SMSCode:        "123456",
		BindPhoneToken: bindPhoneToken,
	}
	w = httptest.Post("/api/login/wechat/qr_code/bind/phone", rb)
	t.Logf("resonse3: %s", w.Body)
	r = responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.NotEmpty(t, body["redirect_url"].(string))
		assert.Contains(t, body["redirect_url"].(string), "http")
		assert.Contains(t, body["redirect_url"].(string), "code")
		assert.Contains(t, body["redirect_url"].(string), "state")
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}
}

func TestWeChatQRCodeLoginBindPhoneExistUser(t *testing.T) {
	var err error
	sceneID := rand.Intn(999999999)
	err = wechatqrcode.SetScene(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetScanOK(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetNeedBindPhone(sceneID, wechatqrcode.WeChatInfo{
		City:       "武汉市",
		Province:   "湖北省",
		Sex:        1,
		Country:    "中国",
		HeadImgURL: "",
		Nickname:   "小白",
		OpenID:     randc.RandStringN(28),
	})
	assert.Nil(t, err)
	reqBody := contexts.WeChatQRCodeLoginReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		SceneId: sceneID,
	}
	w := httptest.Get("/api/login/wechat/qr_code/login", reqBody)
	t.Logf("resonse2: %s", w.Body)
	r := responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	bindPhoneToken := ""
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.False(t, body["auth_ok"].(bool))
		assert.True(t, body["need_bind_phone"].(bool))
		assert.NotEmpty(t, body["bind_phone_token"].(string))
		bindPhoneToken = body["bind_phone_token"].(string)
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}

	// 绑定
	rb := contexts.BindPhoneReq{
		BaseAuthReq: contexts.BaseAuthReq{
			ClientId:    2,
			AuthScene:   1,
			SPMFullCode: "1024.DFASDF234FDAS231.100.101",
			RedirectUrl: "http://www.baidu.com",
			State:       "test",
		},
		// 填一个注册了但是没有绑定微信的手机号
		Phone:          "17852000001",
		SMSCode:        "123456",
		BindPhoneToken: bindPhoneToken,
	}
	w = httptest.Post("/api/login/wechat/qr_code/bind/phone", rb)
	t.Logf("resonse3: %s", w.Body)
	r = responses.Response{}
	err = json.Unmarshal(w.Body.Bytes(), &r)
	assert.Nil(t, err)
	if body, ok := r.Body.(map[string]interface{}); ok {
		assert.NotEmpty(t, body["redirect_url"].(string))
		assert.Contains(t, body["redirect_url"].(string), "http")
		assert.Contains(t, body["redirect_url"].(string), "code")
		assert.Contains(t, body["redirect_url"].(string), "state")
	} else {
		t.Error("响应格式错误", r.Body)
		t.FailNow()
	}
}
