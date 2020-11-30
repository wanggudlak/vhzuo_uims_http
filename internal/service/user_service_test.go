package service_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/db"
	"uims/pkg/encryption"
)

func TestUserService_SetPassword(t *testing.T) {
	var user model.User
	err := db.Def().Where("account = ?", "admin").First(&user).Error
	assert.Nil(t, err)
	t.Logf("user: %v", user)
	switch user.EncryptType {
	case 0:
		// 默认鉴权 / 结算系统鉴权
		assert.True(t, encryption.BcryptCheck("123456", user.Passwd))
	case 1:
		// 任务系统鉴权
		assert.True(t, encryption.DefaultPBKDF2Options.CheckPBKDF2PasswdForVzhuoTaskSYS("123456", user.Passwd, user.Salt))
	default:
		t.Errorf("无效的 encrypt: %d", user.EncryptType)
		t.FailNow()
	}

	err = service.UserService{}.SetPassword(user.ID, "123456")
	assert.Nil(t, err)

	var user2 model.User
	err = db.Def().Where("id = ?", user.ID).First(&user2).Error
	assert.Nil(t, err)
	assert.Equal(t, 0, user2.EncryptType)
	assert.True(t, encryption.BcryptCheck("123456", user2.Passwd))
}

func TestIsExistByPhone(t *testing.T) {
	phone := "13641337591ss"
	exist, err := service.IsExistByPhone(phone)
	assert.Nil(t, err)
	assert.Equal(t, false, exist)
}

func TestIsExistByEmail(t *testing.T) {
	email := "342448932@qq.com"
	exist, err := service.IsExistByEmail(email)
	assert.Nil(t, err)
	assert.Equal(t, true, exist)
}

func TestUserService_UpdateUser(t *testing.T) {
	user := model.User{}
	err := db.Def().Where("account = ?", "zhan").First(&user).Error
	assert.Nil(t, err)
	err = service.GetUserService().UpdateUser(user.ID, 0)
	assert.Nil(t, err)
	newUser := model.User{}
	err = db.Def().Where("id = ?", user.ID).First(&newUser).Error
	assert.Nil(t, err)
	assert.Equal(t, "N", newUser.Status)
	err = db.Def().Model(&newUser).Where("id = ?", user.ID).UpdateColumn("status", "Y").Error
	assert.Nil(t, err)
}
