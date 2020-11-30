package upload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"uims/app"
	"uims/pkg/tool"
)

func PathExt(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// IsNotExistMkDir create a directory if it does not exist
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir create a directory
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// CheckNotExist check if the file exists
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// CheckPermission check if the file has permission
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// GetFileName get file name
func GetFileName(name string) string {
	ext := PathExt(name)
	fileName := strings.TrimSuffix(name, ext)
	//fileName = EncodeMD5(fileName)

	return fileName + ext
}

// GetFileFullPath get full save path
func GetFileFullPath(path string) string {
	return app.StoragePath + "/key/" + path + "/"
}

// CheckFile check if the file exists
func CheckFile(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

func MakeFileDir(path string) (string, error) {
	fileSavedRootPathDir := strings.TrimSuffix(GetFileFullPath(path), "")
	if isExist, _ := tool.IsExistPath(fileSavedRootPathDir); !isExist {
		e := os.MkdirAll(fileSavedRootPathDir, os.ModePerm)
		if e != nil {
			return fileSavedRootPathDir, errors.New("Create dir failed for uploaded file: " + e.Error())
		}
	}
	return fileSavedRootPathDir, nil
}
