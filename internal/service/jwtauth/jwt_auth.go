package jwtauth

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/gjwt"
	"uims/pkg/gredis"
	"uims/pkg/randc"
)

const CodeDataSaveKey = "user_login_code"
const CodeTTL = 300 * time.Second
const (
	// 2 小时有效
	AccessTokenTTL = 7200 * time.Second
	// 30天有效期
	RefreshTokenTTL = 30 * 24 * 60 * 60 * time.Second
)

var FreezeErr = errors.New("用户已冻结")

type UserJwtAuth struct {
	OpenId   string `json:"open_id"`
	ClientId uint   `json:"client_id"`
	Account  string `json:"account"`
	State    string `json:"state"`
}

// 生成 code 码
// 长度为 32
// 5分钟有效期, 储存在 redis 中
// 通过 code 可以反解出需要的数据
func (u *UserJwtAuth) GenerateCode() string {
	code := randc.UUID()
	gredis.Def().Set(u.getRedisKeyName(code), u.JSON(), CodeTTL).Val()
	return code
}

// code 码转换为用户标识
func (u *UserJwtAuth) ParseCode(code string) error {
	var err error
	s, err := gredis.Def().Get(u.getRedisKeyName(code)).Result()
	if err != nil {
		return errors.New("不存在的code")
	}
	if s == "" {
		return errors.New("无效的code")
	}
	err = json.Unmarshal([]byte(s), &u)
	if err != nil {
		log.Error(err)
		return errors.New("解析code失败")
	}
	return nil
}

// 删除code
func (u *UserJwtAuth) RemoveCode(code string) {
	_, _ = gredis.Def().Del(u.getRedisKeyName(code)).Result()
}

func (u *UserJwtAuth) GenerateAccessToken() (string, error) {
	j := AccessClaims{
		OpenId:   u.OpenId,
		ClientId: u.ClientId,
	}
	j.SetAudience("access_token")
	j.SetTTL(AccessTokenTTL)
	j.SetIssue()
	return gjwt.CreateToken(&j)
}

func (u *UserJwtAuth) GenerateRefreshToken() (string, error) {
	j := AccessClaims{
		OpenId:   u.OpenId,
		ClientId: u.ClientId,
	}
	j.SetAudience("refresh_token")
	j.SetTTL(RefreshTokenTTL)
	j.SetIssue()
	return gjwt.CreateToken(&j)
}

func (UserJwtAuth) getRedisKeyName(key string) string {
	return fmt.Sprintf("%s:%s", CodeDataSaveKey, key)
}

func (u *UserJwtAuth) JSON() string {
	b, _ := json.Marshal(u)
	return string(b)
}

// 返回用户是否已经冻结
func (u *UserJwtAuth) IsFreeze() bool {
	user := model.User{}
	db.Def().Select("status").Where("open_id = ?", u.OpenId).First(&user)
	return user.Status != "Y"
}
