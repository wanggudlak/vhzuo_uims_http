package service_test

import (
	"testing"
	"uims/internal/service"
)

func TestUnzip(t *testing.T) {
	srcFile := "/Users/alei/go/src/gitee.com/skysharing/uims-project/uims/storage/app/public/resource/1024.DASFSADFDASFASDF.102.102.zip"
	destFile := "/Users/alei/go/src/gitee.com/skysharing/uims-project/uims/storage/app/public/resource"

	_, err := service.Unzip(srcFile, destFile)
	if err != nil {
		t.Fatalf("Test Unzip failed: %s\n", err.Error())
	}
}
