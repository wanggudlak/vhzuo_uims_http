package service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"uims/internal/controllers/client_controller/requests"
	"uims/internal/model"
	"uims/pkg/db"
)

var (
	ModelContainerCannotEmpty = errors.New("用于存放返回结构的模型对象不能为空")
	QueryWhereCannotEmpty     = errors.New("查询条件不能为空")
)

type ClientSettingService struct {
}

// ExistClientSettingByID checks if there is a tag with the same id
func (ClientSettingService) ExistClientSettingByID(id int) bool {
	var clientSetting model.ClientSetting
	err := db.Def().Select("id").Where("id = ?", id).First(&clientSetting).Error
	if err != nil {
		return false
	}
	return true
}

//  查询客户端设置详细信息
func (ClientSettingService) GetClientSettingByID(id int) (*model.ClientSetting, error) {
	var clientSetting model.ClientSetting
	err := db.Def().Where("id = ?", id).First(&clientSetting).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &clientSetting, nil
}

//  查询客户端设置详细信息
func (ClientSettingService) GetClientSettingByClientID(clientID uint) ([]model.ClientSetting, error) {
	var clientSetting []model.ClientSetting
	err := db.Def().Where("client_id = ?  and isdel = 'N'", clientID).Find(&clientSetting).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return clientSetting, nil
}

// GetClientSettingsByMap 根据业务配置数据类型及SPM全编码查找客户端配置数据
func GetClientSettingsByMap(pClientSetting *model.ClientSetting, where map[string]interface{}, fields []string) error {
	if nil == pClientSetting {
		return ModelContainerCannotEmpty
	}
	if nil == where {
		return QueryWhereCannotEmpty
	}
	err := db.Def().
		Where(where).
		Select(fields).
		First(pClientSetting).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return err
}

// 查询同一客户端是否有重复的Type
func (ClientSettingService) ExistClientByType(id int, Type string) bool {
	var clientSetting model.ClientSetting
	err := db.Def().Select("id").Where("client_id = ? AND type = ? and isdel = 'N'", id, Type).First(&clientSetting).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

// Add Client
func (ClientSettingService) AddClientSetting(this *requests.NewClientSettingRequest) error {
	//频道ID，对于登录业务，频道ID为100；注册业务频道ID为200
	//页面ID，对于登录业务的登录页面ID为101；注册业务的注册页面ID为201
	var (
		channelID = "200"
		pageID    = "201"
	)

	if this.Type == "LGN" {
		channelID, pageID = "100", "101"
	}

	client, e := GetClientService().GetClientByID(this.ClientID)
	if e != nil {
		return e
	}

	// spm编码，由以下组成：client_spm1_code.client_spm2_code.频道ID.页面ID
	spmFullCode := client.ClientSpm1Code + "." + client.ClientSpm2Code + "." + channelID + "." + pageID

	clientSetting := model.ClientSetting{
		BusChannelID: channelID,
		PageID:       pageID,
		ClientID:     uint(this.ClientID),
		//FormFields:       model.NewFieldsMap(this.Fields),  // 磊哥说先不做 留空
		//PageTemplateFile: &this.TemplateFile,
		//PageTemplateFile: this.TemplateFile,
		PageTemplateFile: &model.PageTemplateFile{A: this.TemplateFile, B: ""},
		SpmFullCode:      spmFullCode,
		Type:             this.Type,
	}
	if err := db.Def().Create(&clientSetting).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// 更新客户端设置信息
func (ClientSettingService) UpdateClientSetting(this *requests.ClientSettingRequest, clientID int) error {

	var (
		channelID = "200"
		pageID    = "201"
	)

	if this.Type == "LGN" {
		channelID, pageID = "100", "101"
	}

	client, e := GetClientService().GetClientByID(clientID)
	if e != nil {
		return e
	}

	// spm编码，由以下组成：client_spm1_code.client_spm2_code.频道ID.页面ID
	spmFullCode := client.ClientSpm1Code + "." + client.ClientSpm2Code + "." + channelID + "." + pageID

	clientSetting := model.ClientSetting{
		ID:           uint(this.ID),
		BusChannelID: channelID,
		PageID:       pageID,
		//FormFields:       model.NewFieldsMap(this.Fields),
		//PageTemplateFile: &this.TemplateFile,
		//PageTemplateFile: this.TemplateFile,
		PageTemplateFile: &model.PageTemplateFile{A: this.TemplateFile, B: ""},
		SpmFullCode:      spmFullCode,
		Type:             this.Type,
		CommonModel: &model.CommonModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Def().Model(&model.ClientSetting{}).UpdateColumns(clientSetting).Error; err != nil {
		return err
	}
	return nil
}

// DeleteClientSetting delete a clientSetting
func (ClientSettingService) DeleteClientSetting(id int) error {
	if err := db.Def().Model(&model.ClientSetting{}).Where("id = ?", id).UpdateColumn("isdel", "Y").Error; err != nil {
		return err
	}
	return nil
}

// 所有已经入驻并且已经设置了需要渲染模板页的客户端spm编码等数据
var AllClientsNeedRenderHTML = []model.ClientSetting{}

// GetAllClientsNeedRenderHTML 查询所有已经入驻并且已经设置了需要渲染模板页的客户端spm编码等数据
func GetAllClientsNeedRenderHTML(pClientSettings *[]model.ClientSetting, fields []string) error {
	err := db.Def().
		Where("type IN (?)", []string{"LGN", "REG"}).
		Where("spm_full_code <> ?", "").
		Where("page_template_file <> ?", "").
		Select(fields).
		Find(pClientSettings).
		Error
	if err != nil {
		return err
	}
	return nil
}
