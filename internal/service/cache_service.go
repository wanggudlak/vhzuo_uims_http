package service

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	redissdk "github.com/go-redis/redis/v7"
	"time"
	"uims/internal/model"
	"uims/pkg/const_definition"
	"uims/pkg/gredis"
)

type RedisCacheService struct {
	client *redissdk.Client
}

var RedisDefaultCache = &RedisCacheService{
	client: gredis.Def(),
}

// RedisCacheString 用默认redis驱动缓存字符串型的数据
func (redisCacheService *RedisCacheService) RedisCacheString(key, value string, duration time.Duration) error {
	return redisCacheService.client.Set(key, value, duration).Err()
}

// RedisGetStr 从redis中获取字符串型的缓存值
func (redisCacheService *RedisCacheService) RedisGetStr(key string) string {
	return redisCacheService.client.Get(key).Val()
}

func (redisCacheService *RedisCacheService) RedisDel(key ...string) error {
	return redisCacheService.client.Del(key...).Err()
}

// RedisRememberString 如果缓存中有直接返回这个值，否则重新生成后放入缓存；另外，如果isRefushCache=true，每次都会重新生成一个值，并
// 刷新缓存
func (redisCacheService *RedisCacheService) RedisRememberString(key string, duration time.Duration,
	isRefushCache bool, callback func() interface{}) (string, error) {

	var value string

	if isRefushCache {
		goto GEN_NEW_VALUE
	}

	value = redisCacheService.client.Get(key).Val()
	if len(value) != 0 {
		return value, nil
	}

GEN_NEW_VALUE:
	value = callback().(string)

	err := redisCacheService.client.Set(key, value, duration).Err()

	return value, err
}

func RedisKey(oriK string) string {
	return "uims:" + oriK
}

func MakeSMSRedisCacheKey(phone, captchaID string) string {
	return RedisKey("sms:" + phone + captchaID)
}

//缓存后台登陆用户基本信息
func CacheBackgroundUserInfo(pUserModel *model.User, account string) error {
	cacheKey := const_definition.BACKGROUND_USER_INFO_CACHE_PREFIX + account
	cacheData, _ := json.Marshal(&pUserModel)
	fmt.Println(cacheData)
	err := gredis.Def().Set(
		cacheKey,
		string(cacheData),
		const_definition.BACKGROUND_USER_INFO_EXPIRED*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

//获取后台用户登陆基本信息
func GetBackgroundUserInfoFromCache(account string) (userModel model.User, err error) {
	mapCacheData := gredis.Def().Get(const_definition.BACKGROUND_USER_INFO_CACHE_PREFIX + account).Val()
	cacheData := []byte(mapCacheData)
	_ = json.Unmarshal(cacheData, &userModel)
	//fmt.Printf("u: %+v \n", userModel)

	return userModel, err
}
