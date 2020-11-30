package service_test

import (
	"fmt"
	"testing"
	"uims/internal/service"
)

func TestGetStaticFileTimeChildDir(t *testing.T) {
	r := service.GetStaticFileTimeChildDir("1024.DFASDF234FDAS231.100.101",
		"/resource/1024.DFASDF234FDAS231.100.101/20200630163952/html_template/index.html")
	fmt.Println("r=", r)
}
