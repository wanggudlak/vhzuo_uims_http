package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"uims/internal/controllers/client_controller/requests"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/tool"
	//"uims/tool"
)

type ClientService struct {
}

// GetClientByID checks if there is a tag with the same name
func (ClientService) GetClientByID(id int) (*model.Client, error) {
	var client model.Client
	err := db.Def().Where("id = ?", id).First(&client).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &client, nil
}

// ExistClientByID checks if there is a tag with the same id
func (ClientService) ExistClientByID(id int) bool {
	var client model.Client
	err := db.Def().Select("id").Where("id = ?", id).First(&client).Error
	if err != nil {
		return false
	}
	return true
}

// ExistClientByName checks if there is a tag with the same name
func (ClientService) ExistClientByName(name string) bool {
	var client model.Client
	err := db.Def().Select("id").Where("client_name = ?", name).First(&client).Error
	if err != nil {
		return false
	}
	return true
}

// Add Client
func (ClientService) AddClient(this *requests.ClientNewRequest) error {
	//ForgetTime, err := time.ParseInLocation("2006-01-02 15:04:05", this.ForgetAT, time.Local)
	//if err != nil {
	//	return err
	//}

	var err error
	INTime := time.Now()
	if this.INAT != "" {
		INTime, err = time.ParseInLocation("2006-01-02 15:04:05", this.INAT, time.Local)
		if err != nil {
			return err
		}
	}

	//var ips, urls []byte
	//if this.HostIP != "" {
	//	ips, _ = tool.JSON(this.HostIP)
	//}
	//if this.HostURL != "" {
	//	urls, _ = tool.JSON(this.HostURL)
	//}
	fmt.Println(this.HostIP)
	appID := tool.GenerateRandStrWithMath(16)
	client := model.Client{
		AppId:          appID,
		AppSecret:      "", // 密钥
		ClientType:     this.Type,
		ClientFlagCode: this.FlagCode,
		ClientSpm1Code: this.Spm1Code,
		ClientSpm2Code: appID,
		ClientName:     this.Name,
		ClientHostIp:   this.HostIP,
		ClientHostUrl:  this.HostURL,
		InAt:           INTime,
		//ForgetAt:       ForgetTime,
	}
	if err := db.Def().Create(&client).Error; err != nil {
		return err
	}

	//给客户端创建默认的组织
	org := model.Org{
		ClientID:    client.ID,
		ParentOrgID: 0,
		ClientAppID: client.AppId,
		OrgNameCN:   fmt.Sprintf("%s %s", client.ClientName, "组织"),
	}
	if err := db.Def().Create(&org).Error; err != nil {
		return err
	}

	return nil

}

// GetClients gets a list of tags based on paging and constraints
func (ClientService) GetClients(pageNum int, pageSize int) ([]model.Client, error) {
	var (
		client []model.Client
		err    error
	)
	if pageSize > 0 && pageNum > 0 {
		err = db.Def().Order("id desc").Find(&client).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Def().Order("id desc").Find(&client).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return client, nil
}

// GetClientTotal gets a list of tags based on paging and constraints
func (ClientService) GetClientTotal() (int, error) {
	var count int
	if err := db.Def().Model(&model.Client{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 修改客户端状态  默认N：未授权不可用；Y：已授权可用；F-被禁用
func (ClientService) ChangeClientStatus(this requests.ClientStatusRequest) error {
	if err := db.Def().Model(&model.Client{}).Where("id = ?", this.ID).UpdateColumn("status", this.Status).Error; err != nil {
		return err
	}
	return nil
}

// GetClientId get clientId
func (ClientService) GetClientId(whereMaps interface{}) (clientId uint) {
	var client model.Client

	err := db.Def().
		Select("id").
		Where(whereMaps).
		First(&client).
		Error

	if err == nil {
		clientId = client.ID
	}

	return
}

func (ClientService) UpdateClient(id int, data interface{}) error {
	if err := db.Def().Model(&model.Client{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
