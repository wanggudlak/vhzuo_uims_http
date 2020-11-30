package wechatcontroller

import "uims/pkg/gjwt"

type BindReq struct {
	State string `json:"state" form:"state" comment:"场景" binding:"required"`
	Code  string `json:"code" form:"code" comment:"微信授权码"`
}

type UnbindReq struct {
	Token string `json:"token" form:"token" comment:"解绑使用token" binding:"required"`
}

type BindCache struct {
	ClientId    uint   `json:"client_id"`
	WeChatId    uint   `json:"wechat_id"`
	UserId      uint   `json:"user_id"`
	RedirectURL string `json:"redirect_url"`
	State       string `json:"state"`
}

type UnbindCache struct {
	ClientId    uint   `json:"client_id"`
	WeChatId    uint   `json:"wechat_id"`
	UserId      uint   `json:"user_id"`
	RedirectURL string `json:"redirect_url"`
	State       string `json:"state"`
}

type UnbindClaims struct {
	gjwt.Jwt
	ClientId    uint   `json:"client_id"`
	WeChatId    uint   `json:"wechat_id"`
	UserId      uint   `json:"user_id"`
	RedirectURL string `json:"redirect_url"`
	State       string `json:"state"`
}
