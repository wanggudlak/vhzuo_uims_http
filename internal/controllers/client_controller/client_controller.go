package client_controller

import (
	"encoding/base64"
	"io/ioutil"
	"uims/conf"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/tool"

	//"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
	requests2 "uims/internal/controllers/client_controller/requests"
	responses2 "uims/internal/controllers/responses"
	"uims/internal/service"
	//"uims/conf"
	//"uims/pkg/tool"
)

// @Summary 获取所有入驻客户端系统
// @Produce  json
// @Param page query int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client/list [GET]
func List(c *gin.Context) {

	pageNum := tool.GetPage(c)
	pageSize := conf.AppSetting.PageSize

	result, err := service.GetClientService().GetClients(pageNum, pageSize)

	if err != nil {
		responses2.Error(c, err)
		return
	}
	var clents []interface{}

	for _, client := range result {
		var org model.Org
		//var URL []string
		//if json.Unmarshal(client.ClientHostUrl, &URL) != nil {
		//	responses2.Error(c, errors.New("序列化失败"))
		//	return
		//}
		//
		//for i, v := range URL {
		//	URL[i] = strings.Trim(v, `"`)
		//}

		if client.ClientHostUrl == "null" {
			client.ClientHostUrl = ""
		}
		data := make(map[string]interface{})
		data["id"] = client.ID
		data["app_id"] = client.AppId
		data["app_secret"] = client.AppSecret
		data["type"] = client.ClientType
		data["flag_code"] = client.ClientFlagCode
		data["spm1_code"] = client.ClientSpm1Code
		data["spm2_code"] = client.ClientSpm2Code
		data["name"] = client.ClientName
		data["status"] = client.Status
		data["host_url"] = client.ClientHostUrl
		data["in_at"] = client.InAt.Format("2006-01-02 15:04:05")
		data["forget_at"] = client.ForgetAt.Format("2006-01-02 15:04:05")
		data["created_at"] = client.CreatedAt.Format("2006-01-02 15:04:05")
		data["updated_at"] = client.UpdatedAt.Format("2006-01-02 15:04:05")

		// 获取客户端组织id
		err := db.Def().Where("client_id = ?", client.ID).First(&org).Error
		if err != nil {
			data["org_id"] = 1
		}
		data["org_id"] = org.ID

		clents = append(clents, data)
	}

	total, e := service.GetClientService().GetClientTotal()
	if e != nil {
		responses2.Error(c, e)
	}

	data := make(map[string]interface{})
	data["data"] = clents
	data["total"] = total

	responses2.Success(c, "success", data)
}

// @Summary 获取客户端信息
// @Produce  json
// @Param id query int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client [GET]
func Detail(c *gin.Context) {

	var request requests2.ClientDetailRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	//查询客户端
	client, e := service.GetClientService().GetClientByID(request.ID)
	if e != nil {
		responses2.Error(c, e)
		return
	}
	//var IP, URL []string
	//if json.Unmarshal(client.ClientHostIp, &IP) != nil {
	//	responses2.Error(c, e)
	//	return
	//}
	//
	//if json.Unmarshal(client.ClientHostUrl, &URL) != nil {
	//	responses2.Error(c, e)
	//	return
	//}
	//
	//for i, v := range IP {
	//	IP[i] = strings.Trim(v, `"`)
	//}
	//for i, v := range URL {
	//	URL[i] = strings.Trim(v, `"`)
	//}
	clientSetting, err := service.GetClientSettingService().GetClientSettingByClientID(client.ID)

	if err != nil {
		responses2.Error(c, e)
		return
	}

	data := make(map[string]interface{})

	data["id"] = client.ID
	data["app_id"] = client.AppId
	data["app_secret"] = client.AppSecret
	data["type"] = client.ClientType
	data["flag_code"] = client.ClientFlagCode
	data["spm1_code"] = client.ClientSpm1Code
	data["spm2_code"] = client.ClientSpm2Code
	data["name"] = client.ClientName
	data["status"] = client.Status
	data["host_ip"] = client.ClientHostIp
	data["host_url"] = client.ClientHostUrl
	data["client_pub_key_path"] = client.ClientPubKeyPath
	data["uims_pub_key_path"] = client.UIMSPubKeyPath
	data["uims_pri_key_path"] = client.UIMSPriKeyPath
	data["in_at"] = client.InAt.Format("2006-01-02 15:04:05")
	data["forget_at"] = client.ForgetAt.Format("2006-01-02 15:04:05")
	data["created_at"] = client.CreatedAt.Format("2006-01-02 15:04:05")
	data["updated_at"] = client.UpdatedAt.Format("2006-01-02 15:04:05")

	var settings []interface{}
	for _, setting := range clientSetting {
		set := make(map[string]interface{})
		set["id"] = setting.ID
		set["bus_channel_id"] = setting.BusChannelID
		set["page_id"] = setting.PageID
		set["spm_full_code"] = setting.SpmFullCode
		set["type"] = setting.Type
		set["page_template_file"] = setting.TemplateFile()
		settings = append(settings, set)
	}

	data["setting"] = settings

	responses2.Success(c, "success", data)
}

// @Summary 新增入驻客户端
// @Produce  json
// @Param type body string true "客户端类型，VDK：微桌"
// @Param flag_code body string true "客户端业务系统标识，VDK_CASS：微桌结算系统等；"
// @Param spm1_code body string true "SPM编码中的第一部分，微桌内部系统用1024；外部系统用2048"
// @Param spm2_code body string true "SPM编码中的第二部分"
// @Param name body string true "客户端业务系统名称"
// @Param host_ip body json true "客户端当前使用的IP，多个用json字符串保存"
// @Param host_url body json true "客户端业务系统当前使用的域名，多个用json字符串保存"
// @Param pub_key_path body string true "客户端业务系统的RSA公钥key文件路径"
// @Param in_at body  string false"入驻可以使用的开始时间点，默认为当前时间"
// @Param forget_at body string true "失效时间"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client [post]
func AddClient(c *gin.Context) {
	var request requests2.ClientNewRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	if !(request.Spm1Code == "1024" || request.Spm1Code == "2048") {
		responses2.Failed(c, "SPM编码中的第一部分，微桌内部系统用1024；外部系统用2048", nil)
	}

	//查询客户端名字  不可重复
	if service.GetClientService().ExistClientByName(request.Name) {
		responses2.Failed(c, "name  not exist", nil)
		return
	}

	//创建数据
	err := service.GetClientService().AddClient(&request)
	if err != nil {
		fmt.Println(err)
		responses2.Failed(c, fmt.Sprintf("%s %s", "add client fail", err), nil)
		return
	}

	responses2.Success(c, "success", nil)
}

// @Summary 新增入驻客户端
// @Produce  json
// @Param id body string true "客户端ID"
// @Param type body string false "客户端类型，VDK：微桌"
// @Param host_ip body json false "客户端当前使用的IP，多个用json字符串保存"
// @Param host_url body json false "客户端业务系统当前使用的域名，多个用json字符串保存"
// @Param pub_key_path body string false "客户端业务系统的RSA公钥key文件路径"
// @Param uims_pub_key_path body string false "UIMS系统的RSA公钥文件路径"
// @Param uims_pri_key_path body string false "UIMS系统的RSA私钥文件路径"
// @Param in_at body  string false"入驻可以使用的开始时间点，默认为当前时间"
// @Param forget_at body string false "失效时间"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client [put]
func UpdataClient(c *gin.Context) {
	var request requests2.ClientUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}
	//查询客户端是否存在
	if !service.GetClientService().ExistClientByID(request.ID) {
		responses2.Failed(c, "client  not exist", nil)
		return
	}
	maps := make(map[string]interface{})
	if request.Type != "" {
		maps["client_type"] = request.Type
	}
	if request.FlagCode != "" {
		maps["client_flag_code"] = request.FlagCode
	}
	if request.HostIP != "" {
		//ip, _ := json.Marshal(strings.Split(strings.Trim(request.HostIP, "[]"), ","))
		maps["client_host_ip"] = request.HostIP
	}
	if request.HostURL != "" {
		//url, _ := json.Marshal(strings.Split(strings.Trim(request.HostURL, "[]"), ","))
		maps["client_host_url"] = request.HostURL
	}
	if request.PUBKryPath != "" {
		maps["client_pub_key_path"] = request.PUBKryPath
		b, err := ioutil.ReadFile(request.PUBKryPath)
		if err != nil {
			responses2.Failed(c, "read file error", nil)
		}
		maps["app_secret"] = base64.StdEncoding.EncodeToString(b)
	}
	if request.UIMSPubKeyPath != "" {
		maps["uims_pub_key_path"] = request.UIMSPubKeyPath
	}
	if request.UIMSPriKeyPath != "" {
		maps["uims_pri_key_path"] = request.UIMSPriKeyPath
	}
	if request.INAT != "" {
		inta, err := time.ParseInLocation("2006-01-02 15:04:05", request.INAT, time.Local)
		if err != nil {
			responses2.Failed(c, "INAT  error", nil)
		}
		maps["in_ta"] = inta
	}
	if request.ForgetAT != "" {
		forget, err := time.ParseInLocation("2006-01-02 15:04:05", request.ForgetAT, time.Local)
		if err != nil {
			responses2.Failed(c, "client  not exist", nil)
		}
		maps["in_ta"] = forget
	}
	maps["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

	//更新数据
	if service.GetClientService().UpdateClient(request.ID, maps) != nil {
		responses2.Failed(c, "update client fail", nil)
		return
	}

	responses2.Success(c, "success", nil)
}

// @Summary 修改客户端入驻状态
// @Produce  json
// @Param id body int true "客户端ID"
// @Param status body  string true "默认N：未授权不可用；Y：已授权可用；F-被禁用"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client [post]
func ChangeClientStatus(c *gin.Context) {
	var request requests2.ClientStatusRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	// if status != 'Y'  OR 'F'  Y 授权  F 禁用
	if !(request.Status == "F" || request.Status == "Y" || request.Status == "N") {
		responses2.Failed(c, "status fail", nil)
		return
	}

	if !service.GetClientService().ExistClientByID(request.ID) {
		responses2.Failed(c, "id  not exist", nil)
		return
	}

	//创建数据
	if service.GetClientService().ChangeClientStatus(request) != nil {
		responses2.Failed(c, "add client fail", nil)
		return
	}

	responses2.Success(c, "success", nil)
}

// @Summary 客户端系统设置信息录入
// @Produce  json
// @Param client_id body int true "客户端ID"
// @Param type body  string true "LGN-用于登录的设置；REG-用于注册的设置"
// @Param form_fields body  json true "{'key1':'value','key2':'value2' ...}"
// @Param page_template_file body  json true "{'a':'value','b':'value'}"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client [post]
func AddClientSetting(c *gin.Context) {
	var request requests2.NewClientSettingRequest
	if err := c.ShouldBind(&request); err != nil {
		fmt.Println(err)
		responses2.Error(c, err)
		return
	}

	//查询客户端 是否存在
	if !service.GetClientService().ExistClientByID(request.ClientID) {
		responses2.Failed(c, "client_id  not exist", nil)
		return
	}

	request.Type = strings.ToUpper(request.Type)

	//查询客户端设置信息 是否存在,如果存在 不可以重复设置
	if service.GetClientSettingService().ExistClientByType(request.ClientID, request.Type) {
		responses2.Failed(c, "client type is exist", nil)
		return
	}

	//创建客户端设置信息数据
	err := service.GetClientSettingService().AddClientSetting(&request)
	if err != nil {
		responses2.Failed(c, fmt.Sprintf("%s %s", "add client fail", err), nil)
		return
	}

	responses2.Success(c, "success", nil)
}

// @Summary 客户端系统设置信息更新
// @Produce  json
// @Param client_id body int true "客户端ID"
// @Param type body  string true "LGN-用于登录的设置；REG-用于注册的设置"
// @Param form_fields body  json true "{'key1':'value','key2':'value2' ...}"
// @Param page_template_file body  json true "{'a':'value','b':'value'}"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client [put]
func UpdateClientSetting(c *gin.Context) {
	// 更新和创建 传的请求字段相同
	var request requests2.ClientSettingRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	//查询客户端设置信息 是否存在
	if !service.GetClientSettingService().ExistClientSettingByID(request.ID) {
		responses2.Failed(c, "client_setting_id  not exist", nil)
		return
	}

	// 查询客户端设置信息 -> 取客户端ID
	clientSetting, e := service.GetClientSettingService().GetClientSettingByID(request.ID)
	if e != nil {
		responses2.Error(c, e)
		return
	}
	request.Type = strings.ToUpper(request.Type)

	//查询客户端设置信息 是否存在相同类型,如果不存在 不可以更新
	if !service.GetClientSettingService().ExistClientByType(int(clientSetting.ClientID), request.Type) {
		responses2.Failed(c, fmt.Sprintf("clitnt_setting type= %s not exist", request.Type), nil)
		return
	}

	//更新客户端设置信息数据
	if service.GetClientSettingService().UpdateClientSetting(&request, int(clientSetting.ClientID)) != nil {
		responses2.Failed(c, "add client_setting fail", nil)
		return
	}
	responses2.Success(c, "success", nil)
}

// @Summary 客户端系统设置信息删除
// @Produce  json
// @Param id body int true "客户端设置ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client [DELETE]
func DeleteClientSetting(c *gin.Context) {
	var request requests2.ClientSettingDeleteRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	//查询客户端设置信息 是否存在
	if !service.GetClientSettingService().ExistClientSettingByID(request.ID) {
		responses2.Failed(c, "client_setting_id  not exist", nil)
		return
	}
	// 软删除 客户端设置信息
	if service.GetClientSettingService().DeleteClientSetting(request.ID) != nil {
		responses2.Failed(c, fmt.Sprintf("delete clitnt_setting fail setting_id = %s", strconv.Itoa(request.ID)), nil)
		return
	}

	responses2.Success(c, "success", nil)
}
