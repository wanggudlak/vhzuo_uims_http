package jwtauth_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/internal/model"
	"uims/internal/service/jwtauth"
	"uims/pkg/db"
	"uims/pkg/gjwt"
	"uims/pkg/randc"
)

func TestUserJwtAuth_GenerateCode(t *testing.T) {
	var err error
	admin := model.User{}
	user := model.User{}
	err = db.Def().Where("account = ?", "admin").First(&admin).Error
	assert.Nil(t, err)
	err = db.Def().Where("account = ?", "zhan").First(&user).Error
	assert.Nil(t, err)

	s := jwtauth.UserJwtAuth{
		OpenId:   admin.OpenID,
		ClientId: 1,
		Account:  admin.Account,
		State:    "state",
	}
	assert.NotEmpty(t, s.JSON())
	code := s.GenerateCode()
	t.Logf("admin generate code: %s", code)
	assert.NotEmpty(t, code)
	newS := jwtauth.UserJwtAuth{}
	err = newS.ParseCode(code)
	assert.Nil(t, err)
	assert.Equal(t, s.ClientId, newS.ClientId)
	assert.Equal(t, s.Account, newS.Account)
	assert.Equal(t, s.State, newS.State)

	// test not exist code
	notExistCode := randc.UUID()
	s2 := jwtauth.UserJwtAuth{}
	err = s2.ParseCode(notExistCode)
	assert.NotNil(t, err)
	assert.Equal(t, "不存在的code", err.Error())

	userAuth := jwtauth.UserJwtAuth{
		OpenId:   user.OpenID,
		ClientId: 1,
		Account:  user.Account,
		State:    "state",
	}
	assert.NotEmpty(t, userAuth.JSON())
	userCode := userAuth.GenerateCode()
	t.Logf("user generate code: %s", userCode)
	assert.NotEmpty(t, code)
}

func TestUserJwtAuth_GenerateAccessToken(t *testing.T) {
	var err error
	s := jwtauth.UserJwtAuth{
		OpenId:   randc.UUID(),
		ClientId: 1,
		Account:  "admin",
		State:    "state",
	}
	token, err := s.GenerateAccessToken()
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	var c jwtauth.AccessClaims
	if err := gjwt.Parse(token, &c); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		assert.Equal(t, s.ClientId, c.ClientId)
		assert.Equal(t, s.OpenId, c.OpenId)
	}

	// 无效的token测试
	var c1 jwtauth.AccessClaims
	if err := gjwt.Parse("", &c1); err == nil {
		t.Error("c1 must failed")
		t.FailNow()
	}

	var c2 jwtauth.AccessClaims
	if err := gjwt.Parse("lhlsda", &c2); err == nil {
		t.Error("c2 must failed")
		t.FailNow()
	} else {
		t.Log(err)
	}
}

func TestUserJwtAuth_RemoveCode(t *testing.T) {
	s := jwtauth.UserJwtAuth{}
	code := s.GenerateCode()
	err := s.ParseCode(code)
	assert.Nil(t, err)
	s.RemoveCode(code)
	err = s.ParseCode(code)
	assert.NotNil(t, err)
	assert.Equal(t, "不存在的code", err.Error())
}

func TestUserJwtAuth_IsFreeze(t *testing.T) {
	user := model.User{}
	err := db.Def().First(&user).Error
	assert.Nil(t, err)
	user.Status = "N"
	err = db.Def().Save(&user).Error
	assert.Nil(t, err)
	s := jwtauth.UserJwtAuth{}
	s.OpenId = user.OpenID
	assert.True(t, s.IsFreeze())

	user.Status = "Y"
	err = db.Def().Save(&user).Error
	assert.Nil(t, err)
	assert.False(t, s.IsFreeze())
}
