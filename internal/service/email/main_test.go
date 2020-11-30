package email_test

import (
	"fmt"
	"testing"
	"uims/internal/service/email"
)

func TestAuthToServer_Send(t *testing.T) {
	body := `
		<html>
		<body>
		<h3>
		测试使用Golang发送邮件
		</h3>
		</body>
		</html>
		`
	//body := "测试使用Golang发送邮件"

	c := &email.Context{
		To:       []string{"342448932@qq.com"},
		Subject:  "测试使用Golang发送邮件",
		BodyType: email.HTMlContentType,
		Body:     body,
	}

	err := c.Send()
	if err != nil {
		fmt.Printf("AuthToServer_Send failed: %s\n", err.Error())
	} else {
		fmt.Println("发送邮件成功")
	}
}
