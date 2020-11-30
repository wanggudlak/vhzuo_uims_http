package jwtauth

import (
	"uims/pkg/gjwt"
)

type RefreshClaims struct {
	gjwt.Jwt
	OpenId   string
	ClientId uint
}
