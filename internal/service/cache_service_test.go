package service_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/boot"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/tool"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

//go test -v internal/service/cache_service_test.go -test.run TestCacheBackgroundUserInfo
func TestCacheBackgroundUserInfo(t *testing.T) {
	var user model.User
	var err error
	account := "uims_super_admin"
	err = service.GetUserService().GetUserInfoByAccount(&user, account)
	if err != nil {
		tool.Dump("账号错误")
		return
	}
	_ = service.CacheBackgroundUserInfo(&user, user.Account)

	tool.Dump("缓存成功")
}

//go test -v internal/service/cache_service_test.go -test.run TestGetBackgroundUserInfoFromCache
func TestGetBackgroundUserInfoFromCache(t *testing.T) {
	var err error
	account := "uims_super_admin"
	cacheData, err := service.GetBackgroundUserInfoFromCache(account)
	if err != nil {
		tool.Dump(err)
	}
	tool.Dump(cacheData)
}

func TestRedisCacheService_RedisCacheString(t *testing.T) {
	err := service.RedisDefaultCache.RedisCacheString("test", "test", 0)
	assert.Nil(t, err)
	query := service.RedisDefaultCache.RedisGetStr("test")
	assert.Equal(t, "test", query)
}
