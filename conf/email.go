package conf

import "uims/pkg/env"

type Email struct {
	Handler      string
	Driver       string
	Host         string
	Port         int
	User         string
	Passwd       string
	From         string
	FromNickname string
}

func NewEmailConf() *Email {
	return &Email{
		Handler:      env.DefGetStr("MAIL_HANDLER", "net/smtp"), // 也可以用gomail
		Driver:       env.DefGetStr("MAIL_DRIVER", "smtp"),
		Host:         env.DefGetStr("MAIL_HOST", "smtp.exmail.qq.com"),
		Port:         env.DefaultGetInt("MAIL_PORT", 465),
		User:         env.DefGetStr("MAIL_USERNAME", "pay@vzhuo.com"),
		Passwd:       env.DefGetStr("MAIL_PASSWORD", "zxcvbnmZXCVBNM20190308"),
		From:         env.DefGetStr("MAIL_FROM_ADDRESS", "pay@vzhuo.com"),
		FromNickname: env.DefGetStr("MAIL_FROM_NICK", "vzhuo"),
	}
}
