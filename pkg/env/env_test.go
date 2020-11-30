package env_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"uims/boot"
	"uims/pkg/env"
	"uims/pkg/tool"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

func TestDefaultGet(t *testing.T) {
	assert.Equal(t, os.Getenv("APP_NAME"), env.DefaultGet("APP_NAME", ""))
}

func TestSetKeyStringV(t *testing.T) {
	err := env.SetKeyStringV("APP_KEY", "base64:QgvICXZ5WHDgL6GSuBN7RTN/QvVt1Z9lxd0GGPcVvhM=aleijuzixiaodou")
	tool.Dump(err)
	assert.Equal(t, nil, err)
}
