package service

import (
	"github.com/pkg/errors"
	"time"
	"uims/pkg/gredis"
	"uims/pkg/randc"
)

const CSRF_FLAG = "csrf"

var (
	CSRFliveTime               = 5 * time.Minute // CSRF TOKEN 存活时间
	GenerateCSRFTokenFailedErr = errors.New("Generate CSRF token failed")
	CSRFTokenHadExpiredErr     = errors.New("页面数据已失效，请重新刷新页面")
	CSRFTokenVerifiedFailedErr = errors.New("当前页面可能已被伪造或篡改，为了安全起见，请重新刷新页面")
	CSRFTokenParsedFailed      = errors.New("CSRF-TOKEN解析失败")
)

func GenerateCSRFToken() (string, error) {
	token := randc.UUID()
	key := makeCSRFCacheKey(token)
	err := RedisDefaultCache.RedisCacheString(key, "1", 5*time.Minute)
	if err != nil {
		return "", errors.Wrap(err, "储存 CSRF token 失败")
	}
	return token, nil
}

func makeCSRFCacheKey(token string) string {
	return RedisKey(CSRF_FLAG + ":" + token)
}

func VerifyCSRFToken(token string) bool {
	key := makeCSRFCacheKey(token)
	return gredis.Def().Exists(key).Val() == 1
}
