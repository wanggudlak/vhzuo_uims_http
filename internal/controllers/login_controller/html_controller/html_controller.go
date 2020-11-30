package html_controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"uims/app"
	"uims/command/commands/version"
	"uims/internal/controllers/login_controller/contexts"
	resp "uims/internal/controllers/responses"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/wechat"
	"uims/pkg/tool"
)

// RenderHTML 渲染HTML页面
func RenderHTML(c *gin.Context) {
	var req contexts.LoginHTMLrenderRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.BadReq(c, err)
		return
	}

	uri := c.Request.URL.Path
	t := service.SplitURI(uri, "/")
	req.SPM = t[len(t)-1]
	pSpm := service.ParseSPMstring(req.SPM)

	// 通过spm编码去查找当前是哪个业务系统
	setting := &model.ClientSetting{}
	where := map[string]interface{}{
		"spm_full_code": pSpm.FullCode,
	}
	err := service.GetClientSettingsByMap(setting, where,
		[]string{"client_id", "page_template_file", "spm_full_code"})
	if nil != err {
		resp.Error(c, err)
		return
	}

	timeChildDir := service.GetStaticFileTimeChildDir(setting.SpmFullCode, setting.TemplateFile())
	staticFilePathPrefix := "/" + req.SPM + "/" + timeChildDir + "/static"

	htmlTemplateFilePathPrefix := app.GetStoragePath("/app/public")
	tmplFilePath := filepath.Join(htmlTemplateFilePathPrefix, setting.TemplateFile())

	isExistThisFile, err := tool.IsExistPath(tmplFilePath)
	if err != nil {
		resp.Error(c, err)
		return
	}
	if !isExistThisFile {
		resp.Error(c, errors.New("未能获取到客户端系统登录模板页"))
		return
	}

	router := app.GetEngineRouter()
	service.CustomizateHTMLtemplateParamDelimiter(router)
	service.LoadHTMLtemplateFiles(router, tmplFilePath)

	csrfToken, err := service.GenerateCSRFToken()
	if err != nil {
		resp.Error(c, err)
		return
	}

	c.HTML(http.StatusOK, filepath.Base(tmplFilePath), gin.H{
		"_token":                 csrfToken,
		"client_id":              setting.ClientID,
		"spm_full_code":          setting.SpmFullCode,
		"static_filepath_prefix": staticFilePathPrefix,
		"redirect_url":           req.RedirectURL,
		"state":                  req.State,
		"wechat_app_id":          getWeChatAppId(int(setting.ClientID)),
		"v":                      version.VersionSN,
	})

	return
}

func getWeChatAppId(clientId int) string {
	var err error
	c, err := wechat.GetConfig(clientId)
	if err != nil {
		return ""
	}
	return c.AppId
}
