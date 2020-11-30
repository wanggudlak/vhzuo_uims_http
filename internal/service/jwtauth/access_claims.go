package jwtauth

import (
	"uims/pkg/gjwt"
)

type AccessClaims struct {
	gjwt.Jwt
	// 用户 uuid
	OpenId   string
	ClientId uint
}
