package wechatserver

import (
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/oauth2"
)

var clients map[string]*oauth2.Client

type Config struct {
	AppId    string
	Secret   string
	WeChatId uint
	ClientId uint
}

// 获取微信客户端
func Cli(c Config) *oauth2.Client {
	if cli, ok := clients[c.AppId]; ok {
		return cli
	} else {
		return Create(c)
	}
}

func Create(c Config) *oauth2.Client {
	return &oauth2.Client{
		Endpoint: mpoauth2.NewEndpoint(c.AppId, c.Secret),
	}
}
