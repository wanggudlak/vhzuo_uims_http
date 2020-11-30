package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
	"uims/app"
	"uims/conf"
	"uims/internal/model"
	"uims/pkg/tool"
)

const (
	LEFT_HTML_TMPL_PARAM_DELIMITER  = "{[{"
	RIGHT_HTML_TMPL_PARAM_DELIMITER = "}]}"
)

const DIR_SEPARATOR = "/"

// CustomizateHTMLtemplateParamDelimiter 设置HTML模板参数的定界符
func CustomizateHTMLtemplateParamDelimiter(router *gin.Engine) {
	router.Delims(LEFT_HTML_TMPL_PARAM_DELIMITER, RIGHT_HTML_TMPL_PARAM_DELIMITER)
}

func LoadHTMLtemplateFiles(router *gin.Engine, tmplFiles ...string) {
	router.LoadHTMLFiles(tmplFiles...)
}

func MakeHTMLtmplateFileDir(childDir string) (string, string, error) {
	fileSavedRootPathDir := strings.TrimSuffix(app.StoragePath, DIR_SEPARATOR) +
		DIR_SEPARATOR +
		strings.TrimPrefix(conf.Filesystems.Disks.Local.Root, DIR_SEPARATOR) +
		"/resource" + "/" + childDir

	if isExist, _ := tool.IsExistPath(fileSavedRootPathDir); !isExist {
		e := os.Mkdir(fileSavedRootPathDir, 0777)
		if e != nil {
			return fileSavedRootPathDir, "", errors.New("Create dir failed for uploaded file: " + e.Error())
		}
	}

	fileSavedRelativePathDir := "/resource" + "/" + childDir

	return fileSavedRootPathDir, fileSavedRelativePathDir, nil
}

func GetStaticFileTimeChildDir(spmFullCode, htmlSavedRelativePath string) string {
	timeChildDir := strings.Split(strings.TrimPrefix(strings.SplitAfter(htmlSavedRelativePath, spmFullCode)[1], "/"), "/")[0]
	return timeChildDir
}

func GetStaticFileRouteURIPrefix() string {
	return DIR_SEPARATOR + "html"
}

func GetStaticFileRouteURI(spm string) string {
	return DIR_SEPARATOR + spm
}

func GetStaticFileStoragePath(staticRouteURI string) string {
	return app.GetStoragePath("app/public/resource" + staticRouteURI)
}

func GetStaticFileRouteSuffix(pClientSetting *model.ClientSetting) string {
	if nil == pClientSetting {
		return ""
	}

	switch pClientSetting.Type {
	case "LGN":
		return DIR_SEPARATOR + "login" + DIR_SEPARATOR + "render" + DIR_SEPARATOR + pClientSetting.SpmFullCode
	case "REG":
		return DIR_SEPARATOR + "register" + DIR_SEPARATOR + "render" + DIR_SEPARATOR + pClientSetting.SpmFullCode
	default:
		return ""
	}
}
