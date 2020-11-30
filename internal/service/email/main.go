package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"net/smtp"
	"strings"
	"uims/conf"
)

const (
	HandlerNetSMTP = "net/smtp"
	HandlerGoMail  = "gomail"

	HTMlContentType = "html"
	TextContentType = "text"
)

type Context struct {
	To       []string // 收件人邮箱
	Subject  string   // 邮件主题
	BodyType string   // 邮件内容类型，html 、text
	Body     string   // 邮件内容
}

func (c *Context) Send() error {
	cf := conf.EmailConf
	switch cf.Handler {
	default:
		fallthrough
	case HandlerGoMail:
		return c.SendByGoMail(cf)
	case HandlerNetSMTP:
		return c.SendByNetHttp(cf)
	}
}

func (c *Context) SendByGoMail(cf *conf.Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", cf.FromNickname+"<"+cf.From+">")
	m.SetHeader("To", c.To...)
	m.SetHeader("Subject", c.Subject)
	m.SetBody("text/html", c.Body)
	d := gomail.NewDialer(cf.Host, cf.Port, cf.User, cf.Passwd)
	return d.DialAndSend(m)
}

func (c *Context) SendByNetHttp(cf *conf.Email) error {
	plainAuth := smtp.PlainAuth("", cf.User, cf.Passwd, cf.Host)
	contentType := GetContentType(c.BodyType)
	msg := []byte("To: " + strings.Join(c.To, ",") +
		"\r\nForm: " + "微桌" + "<" + cf.From + ">" +
		"\r\nSubject: " + c.Subject +
		"\r\n" + contentType +
		"\r\n\r\n" +
		c.Body)

	fmt.Println(c.To)
	fmt.Println(string(msg))

	return smtp.SendMail(fmt.Sprintf("%s:%d", cf.Host, cf.Port), plainAuth, cf.From, c.To, msg)
}

func GetContentType(bodyType string) string {
	switch bodyType {
	case HTMlContentType:
		return "Content-Type: text/" + bodyType + "; charset=UTF-8"
	default:
		return "Content-Type: text/plain" + "; charset=UTF-8"
	}
}
