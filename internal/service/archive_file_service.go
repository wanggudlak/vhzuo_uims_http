package service

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Unzip 解压zip包
// zipFile -- 被解压的zip包全路径
// destDir -- 被解压到的目标目录
func Unzip(zipFile, destDir string) (string, error) {
	resultDirName := ""
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer zipReader.Close()

	for i, f := range zipReader.File {
		if 0 == i {
			resultDirName = f.Name
		}
		fpath := filepath.Join(destDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return resultDirName, err
			}

			inFile, err := f.Open()
			if err != nil {
				return resultDirName, err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return resultDirName, err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return resultDirName, err
			}
		}
	}
	return resultDirName, nil
}
