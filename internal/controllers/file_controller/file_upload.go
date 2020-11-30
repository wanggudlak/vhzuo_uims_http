package file_controller

import (
	"os"
	"time"
	"uims/internal/controllers/file_controller/requests"
	responses2 "uims/internal/controllers/file_controller/responses"
	responses3 "uims/internal/controllers/responses"
	"uims/internal/service"
	"uims/pkg/upload"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"uims/conf"
)

const (
	DIR_SEPARATOR               = "/"
	ALLOW_UPLOADED_TMPLFILE_EXT = ".zip" // 允许上传的HTML模板文件的后缀
)

func SetMaxMultipleMemory(router *gin.Engine) {
	router.MaxMultipartMemory = conf.MAX_MULTIPART_MEMORY
}

// @Summary 上传登录HTML模板页面
// @Produce  json
// @Param Form
// @Success 200 {object} fileresp.FileUploadedResponse
// @Router /api/upload/htmltmpl [POST]
// UploadClientHTMLtmpl
func UploadClientHTMLtmpl(c *gin.Context) {
	var fileUploadRequest requests.FileUploadRequest
	if err := c.Bind(&fileUploadRequest); err != nil {
		responses3.BadReq(c, err)
		return
	}
	log.WithFields(log.Fields{
		"fileUploadRequest": fileUploadRequest,
	}).Info("记录文件上传的请求")

	// 建立一个保存上传的文件的目录
	// storage/app/public/resource/<spmfull>
	fileSavedRootPathDir, fileSavedRelativePathDir, err := service.MakeHTMLtmplateFileDir(fileUploadRequest.SPMfullCode)
	if err != nil {
		responses3.Failed(c, err.Error(), nil)
		return
	}
	// 上传的文件原名
	originUploadedFilename := filepath.Base(fileUploadRequest.UploadedFile.Filename)
	if len(originUploadedFilename) == 0 {
		responses3.Failed(c, "upload file failed", nil)
		return
	}
	extOriginFile := filepath.Ext(originUploadedFilename)
	if extOriginFile != ALLOW_UPLOADED_TMPLFILE_EXT {
		responses3.Failed(c, "上传的模板文件需要打包成zip压缩文件", nil)
		return
	}

	// 文件保存到正确位置时应该具有的全路径
	uploadedFileWillSavedFullPath := fileSavedRootPathDir + DIR_SEPARATOR + fileUploadRequest.SPMfullCode + extOriginFile
	if err := c.SaveUploadedFile(fileUploadRequest.UploadedFile, uploadedFileWillSavedFullPath); err != nil {
		responses3.Error(c, err)
		return
	}
	// 若上传成功，将zip文件解压缩，解压后保存到同名文件夹中
	resultDirName, err2 := service.Unzip(uploadedFileWillSavedFullPath, fileSavedRootPathDir)
	if err2 != nil {
		responses3.Error(c, err2)
		return
	}
	_ = os.Remove(uploadedFileWillSavedFullPath)
	// 将解压出的文件夹重命名
	log.Info("resultDirName=", resultDirName)
	originResultDirName := filepath.Join(fileSavedRootPathDir, resultDirName)
	renameChildDirName := time.Now().Format("20060102150405")
	renameResultDirName := filepath.Join(fileSavedRootPathDir, renameChildDirName)
	err = os.Rename(originResultDirName, renameResultDirName)
	if err != nil {
		responses3.Error(c, err2)
		return
	}

	fileSavedAbsolutePath := filepath.Join(renameResultDirName, "html_template", "index.html")
	fileSavedRelativePath := filepath.Join(fileSavedRelativePathDir, renameChildDirName, "html_template", "index.html")

	fileUploadedResponse := responses2.FileUploadedResponse{}
	fileUploadedResponse.FileSavedRelativePath = fileSavedRelativePath
	fileUploadedResponse.FileSavedAbsolutePath = fileSavedAbsolutePath

	responses3.Success(c, "success", fileUploadedResponse)
	return
}

// @Summary Import Image
// @Produce  json
// @Param image formData file true "File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func UploadFile(c *gin.Context) {

	var fileUploadRequest requests.UploadRequest
	if err := c.Bind(&fileUploadRequest); err != nil {
		responses3.BadReq(c, err)
		return
	}
	log.WithFields(log.Fields{
		"fileUploadRequest": fileUploadRequest,
	}).Info("记录文件上传的请求")

	if !(fileUploadRequest.Key == "client_pub_key_path" ||
		fileUploadRequest.Key == "uims_pub_key_path" ||
		fileUploadRequest.Key == "uims_pri_key_path") {
		responses3.Failed(c, "参数不正确", nil)
		return
	}

	//查询客户端
	client, e := service.GetClientService().GetClientByID(fileUploadRequest.ClientID)
	if e != nil {
		responses3.Error(c, e)
		return
	}

	// opt/data/appid/abc.key
	folder := client.AppId
	// 文件名 fileUploadRequest.Key + 后缀
	fileName := fileUploadRequest.Key + upload.PathExt(fileUploadRequest.UploadedFile.Filename)
	fileRealName := upload.GetFileName(fileUploadRequest.UploadedFile.Filename) // 文件真实姓名,前端展示使用
	fullPath, err1 := upload.MakeFileDir(folder)
	if err1 != nil {
		responses3.Error(c, err1)
		return
	}

	src := fullPath + fileName

	err := upload.CheckFile(fullPath)
	if err != nil {
		responses3.Error(c, err)
		return
	}

	if err := c.SaveUploadedFile(fileUploadRequest.UploadedFile, src); err != nil {
		responses3.Error(c, err)
		return
	}

	responses3.Success(c, "success", map[string]string{
		"file_path": src,
		"key":       fileUploadRequest.Key,
		"file_name": fileRealName,
	})
}
