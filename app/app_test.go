package app_test

import (
	"testing"
	"uims/app"
	"uims/boot"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

func TestPath(t *testing.T) {
	t.Logf("Storage Path: %s", app.GetStoragePath(""))
}
