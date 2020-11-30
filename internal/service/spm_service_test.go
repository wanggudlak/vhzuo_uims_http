package service_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/internal/service"
	"uims/pkg/tool"
)

func TestParseSPMstring(t *testing.T) {
	pSpm := service.ParseSPMstring("1024.DFASDF234FDAS2314243214.100.101")
	tool.Dump(*pSpm)
	assert.Equal(t, service.SPMcode{
		FullCode: "1024.DFASDF234FDAS2314243214.100.101",
		Code1:    "1024",
		Code2:    "DFASDF234FDAS2314243214",
		Code3:    "100",
		Code4:    "101",
	}, *pSpm)
}
