package service_test

import (
	"fmt"
	"testing"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/tool"
)

func TestGetAllClientsNeedRenderHTML(t *testing.T) {
	clientSettings := []model.ClientSetting{}
	err := service.GetAllClientsNeedRenderHTML(&clientSettings, []string{"*"})
	if err != nil {
		t.Fatalf("GetAllClientsNeedRenderHTML error: <%s>\n", err.Error())
	}
	fmt.Printf("%v\n", clientSettings[0])
	fmt.Printf("%v\n", clientSettings[1])
	tool.Dump(clientSettings)
}
