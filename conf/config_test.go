package conf_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"uims/boot"
	"uims/conf"
)

func TestVal(t *testing.T) {
	boot.SetInTest()
	assert.NotEmpty(t, conf.Name)
	assert.NotEmpty(t, os.Getenv("APP_NAME"))
	assert.Equal(t, conf.Name, os.Getenv("APP_NAME"))
	assert.Equal(t, conf.APPKey, os.Getenv("APP_KEY"))
}
