package oauthcontroller_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/boot"
	"uims/internal/model"
	"uims/internal/service/jwtauth"
	"uims/internal/thriftcontroller/oauthcontroller"
	"uims/pkg/db"
	"uims/pkg/gjwt"
	"uims/pkg/randc"
	"uims/pkg/thrift/common"
	thriftserver "uims/pkg/thrift/server"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

func TestAccessToken(t *testing.T) {
	var err error
	jwtAuth := jwtauth.UserJwtAuth{}
	jwtAuth.OpenId = randc.UUID()
	jwtAuth.ClientId = 1
	code := jwtAuth.GenerateCode()
	c := thriftserver.NewContext()
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: oauthcontroller.AccessTokenReq{
			Code:      code,
			GrantType: "authorization_code",
		},
	}.JSON())
	assert.Nil(t, err)
	oauthcontroller.AccessToken(c)
	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Data.Status)
	assert.NotNil(t, c.Response.Data.Content)
	assert.NotNil(t, c.Response.Data.Content.(oauthcontroller.AccessTokenResp))
	assert.Equal(t, jwtAuth.OpenId, c.Response.Data.Content.(oauthcontroller.AccessTokenResp).OpenId)
	assert.NotEmpty(t, c.Response.Data.Content.(oauthcontroller.AccessTokenResp).AccessToken)
	assert.NotEmpty(t, c.Response.Data.Content.(oauthcontroller.AccessTokenResp).RefreshToken)
	assert.NotEmpty(t, c.Response.Data.Content.(oauthcontroller.AccessTokenResp).ExpiresIn)
	t.Logf("%+v", c.Response)
	t.Logf("%+v", c.Response.Data.Content.(oauthcontroller.AccessTokenResp))

	// 解析token

	var claims jwtauth.AccessClaims
	if err := gjwt.Parse(c.Response.Data.Content.(oauthcontroller.AccessTokenResp).AccessToken, &claims); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		assert.Equal(t, jwtAuth.OpenId, claims.OpenId)
		assert.Equal(t, jwtAuth.ClientId, claims.ClientId)
	}

	var refreshClaims jwtauth.RefreshClaims
	if err := gjwt.Parse(c.Response.Data.Content.(oauthcontroller.AccessTokenResp).RefreshToken, &refreshClaims); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		assert.Equal(t, jwtAuth.OpenId, claims.OpenId)
		assert.Equal(t, jwtAuth.ClientId, claims.ClientId)
	}

	// code 应该被删除了
	err = jwtAuth.ParseCode(code)
	assert.NotNil(t, err)
	assert.Equal(t, "不存在的code", err.Error())
}

func TestRefreshToken(t *testing.T) {
	auth := jwtauth.UserJwtAuth{
		ClientId: 1,
		OpenId:   randc.UUID(),
	}

	refreshToken, err := auth.GenerateRefreshToken()
	assert.Nil(t, err)

	c := thriftserver.NewContext()
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: oauthcontroller.RefreshTokenReq{
			RefreshToken: refreshToken,
		},
	}.JSON())
	assert.Nil(t, err)
	oauthcontroller.RefreshToken(c)
	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Data.Status)
	assert.NotNil(t, c.Response.Data.Content)
	assert.NotNil(t, c.Response.Data.Content.(oauthcontroller.AccessTokenResp))
	assert.Equal(t, auth.OpenId, c.Response.Data.Content.(oauthcontroller.AccessTokenResp).OpenId)
	assert.NotEmpty(t, c.Response.Data.Content.(oauthcontroller.AccessTokenResp).AccessToken)
	assert.NotEmpty(t, c.Response.Data.Content.(oauthcontroller.AccessTokenResp).RefreshToken)
	assert.NotEmpty(t, c.Response.Data.Content.(oauthcontroller.AccessTokenResp).ExpiresIn)
}

func TestUserInfo(t *testing.T) {
	var openIds []string
	err := db.Def().Model(&model.User{}).Where("open_id != ?", "").Pluck("open_id", &openIds).Error
	assert.Nil(t, err)
	assert.True(t, len(openIds) > 0)
	auth := jwtauth.UserJwtAuth{
		ClientId: 1,
		OpenId:   openIds[0],
	}
	accessToken, err := auth.GenerateAccessToken()
	assert.Nil(t, err)
	c := thriftserver.NewContext()
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: oauthcontroller.UserInfoReq{
			AccessToken: accessToken,
		},
	}.JSON())
	oauthcontroller.UserInfo(c)
	assert.Nil(t, err)
	t.Logf("response: %+v", c.Response)
	t.Logf("response data: %+v", c.Response.Data)

	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	if resp, ok := c.Response.Data.Content.(oauthcontroller.UserInfoResp); !ok {
		t.Errorf("%+v type not is %s", c.Response.Data.Content, "oauthcontroller.UserInfoResp")
		assert.NotEmpty(t, resp.OpenId)
		t.FailNow()
	} else {
		assert.Equal(t, auth.OpenId, resp.OpenId)
	}
	t.Logf("response thrift: %+v", c.Response.ConvertThriftResp())
}

func TestUserInfoWeChat(t *testing.T) {
	var userIDs []uint
	err := db.Def().Model(&model.UserWeChat{}).Pluck("user_id", &userIDs).Error
	for _, userId := range userIDs {
		var user model.User
		err = db.Def().Where("id = ?", userId).Select("open_id").First(&user).Error
		assert.Nil(t, err)
		assert.NotEmpty(t, user.OpenID)
		auth := jwtauth.UserJwtAuth{
			ClientId: 1,
			OpenId:   user.OpenID,
		}
		accessToken, err := auth.GenerateAccessToken()
		assert.Nil(t, err)
		c := thriftserver.NewContext()
		err = c.ParseRequest(thriftserver.BaseRequest{
			Method: "test",
			Params: oauthcontroller.UserInfoReq{
				AccessToken: accessToken,
			},
		}.JSON())
		oauthcontroller.UserInfo(c)
		assert.Nil(t, err)
		assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
		assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
		assert.Equal(t, common.STATUS_SUCCESS, c.Response.Data.Status)
		if resp, ok := c.Response.Data.Content.(oauthcontroller.UserInfoResp); !ok {
			t.Errorf("%+v type not is %s", c.Response.Data.Content, "oauthcontroller.UserInfoResp")
			t.FailNow()
		} else {
			assert.Equal(t, auth.OpenId, resp.OpenId)
			assert.NotEmpty(t, resp.Phone)
			assert.NotEmpty(t, resp.UserCode)
			assert.True(t, len(resp.WeChats) > 0)
			for _, wechat := range resp.WeChats {
				assert.NotEmpty(t, wechat.OpenID)
				assert.NotEmpty(t, wechat.UUID)
			}
		}
	}
}

func TestGetBindWeChatURL(t *testing.T) {
	var err error
	var user model.User
	err = db.Def().First(&user).Error
	assert.Nil(t, err)
	claims := jwtauth.AccessClaims{
		OpenId:   user.OpenID,
		ClientId: 1,
	}
	token, err := gjwt.CreateToken(&claims)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	c := thriftserver.NewContext()
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: oauthcontroller.GetWeChatBindURLReq{
			AccessToken: token,
			WeChatUUID:  "463c1a8bdc6c4f07a94290b524cd559c",
			RedirectURL: "http://web.vzhuo.com/loading",
		},
	}.JSON())
	assert.Nil(t, err)
	oauthcontroller.GetBindWeChatURL(c)
	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	if resp, ok := c.Response.Data.Content.(oauthcontroller.GetWeChatBindURLResp); !ok {
		t.Errorf("%+v type not is %s", c.Response.Data.Content, "oauthcontroller.GetWeChatBindURLResp")
		t.FailNow()
	} else {
		t.Log(resp.URL)
		assert.NotEmpty(t, resp.URL)
	}
}

func TestGetUnbindWeChatURL(t *testing.T) {
	var err error
	var user model.User
	var userWeChat model.UserWeChat
	var wechat model.WeChat
	var wechatUUID = "463c1a8bdc6c4f07a94290b524cd559c"
	err = db.Def().Where("uuid = ?", wechatUUID).First(&wechat).Error
	assert.Nil(t, err)
	err = db.Def().Where("wechat_id = ?", wechat.ID).First(&userWeChat).Error
	assert.Nil(t, err)
	err = db.Def().Where("id = ?", userWeChat.UserId).First(&user).Error
	assert.Nil(t, err)
	claims := jwtauth.AccessClaims{
		OpenId:   user.OpenID,
		ClientId: 1,
	}
	token, err := gjwt.CreateToken(&claims)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	c := thriftserver.NewContext()
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: oauthcontroller.GetWeChaUnbindURLReq{
			AccessToken: token,
			WeChatUUID:  wechatUUID,
			RedirectURL: "http://web.vzhuo.com/loading",
		},
	}.JSON())
	assert.Nil(t, err)
	oauthcontroller.GetUnbindWeChatURL(c)
	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	if resp, ok := c.Response.Data.Content.(oauthcontroller.GetWeChatUnbindURLResp); !ok {
		t.Logf("%+v", c.Response)
		t.Errorf("%+v type not is %s", c.Response.Data.Content, "oauthcontroller.GetWeChatUnbindURLResp")
		t.FailNow()
	} else {
		t.Log(resp.URL)
		assert.NotEmpty(t, resp.URL)
	}
}

func TestIsBindWeChat(t *testing.T) {
	var openIds []string
	err := db.Def().Model(&model.User{}).Where("open_id != ?", "").Pluck("open_id", &openIds).Error
	assert.Nil(t, err)
	assert.True(t, len(openIds) > 0)
	auth := jwtauth.UserJwtAuth{
		ClientId: 1,
		OpenId:   openIds[0],
	}
	accessToken, err := auth.GenerateAccessToken()
	assert.Nil(t, err)
	c := thriftserver.NewContext()
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: oauthcontroller.IsBindWeChatReq{
			AccessToken: accessToken,
			WeChatUUID:  "463c1a8bdc6c4f07a94290b524cd559c",
		},
	}.JSON())
	oauthcontroller.IsBindWeChat(c)
	assert.Nil(t, err)
	t.Logf("response: %+v", c.Response)
	t.Logf("response data: %+v", c.Response.Data)

	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	if resp, ok := c.Response.Data.Content.(oauthcontroller.IsBindWeChatResp); !ok {
		t.Log(resp)
		t.Errorf("%+v type not is %s", c.Response.Data.Content, "oauthcontroller.IsBindWeChatResp")
		t.FailNow()
	} else {
		assert.False(t, resp.IsBind)
	}
	t.Logf("response thrift: %+v", c.Response.ConvertThriftResp())
}

func TestQueryUserWeChat(t *testing.T) {
	var err error
	var weChats = []oauthcontroller.UserInfoWeChat{}
	err = db.Def().
		Table("uims_user_wechat as uw").
		Select([]string{"uw.nickname", "uw.wechat_open_id as open_id", "uw.avatar", "w.uuid"}).
		Joins("left join uims_wechat as w on w.id = uw.wechat_id").
		Where("uw.user_id = ?", 75).
		Scan(&weChats).Error
	assert.Nil(t, err)
	t.Log(weChats)
}
