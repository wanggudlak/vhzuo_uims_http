package gjwt

import (
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"time"
	"uims/conf"
)

// 可以便捷的定义 Claims, 生成 Token 与 解析为 Claims
// 使用方式参照 gjwt_test.go
type JWTer interface {
	jwtgo.Claims
	SetTTL(t time.Duration)
	SetAudience(audience string)
	SetIssue()
}

func CreateToken(j JWTer) (string, error) {
	t := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, j)
	return t.SignedString([]byte(conf.APPKey))
}

func Parse(token string, jwt JWTer) error {
	if token == "" {
		return fmt.Errorf("token 不允许为空")
	}
	t, err := jwtgo.ParseWithClaims(token, jwt, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(conf.APPKey), nil
	})
	if err != nil {
		return err
	}
	if !t.Valid {
		return jwt.Valid()
	}
	return nil
}

// 用来给其它结构体组合
type Jwt struct {
	jwtgo.StandardClaims
}

func (j *Jwt) SetTTL(t time.Duration) {
	now := time.Now()
	if t != 0 {
		expireTime := now.Add(t)
		j.ExpiresAt = expireTime.Unix()
	}
}

func (j *Jwt) SetAudience(audience string) {
	j.Audience = audience
}

func (j *Jwt) SetIssue() {
	now := time.Now()
	j.IssuedAt = now.Unix()
	j.Issuer = conf.Name
}
