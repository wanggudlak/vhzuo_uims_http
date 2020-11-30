package service

import (
	"github.com/jinzhu/gorm"
	"uims/internal/model"
	"uims/pkg/db"
)

type OrgService struct {
}

// 查询ORG是否存在
func (OrgService) ExistOrgByID(id int) bool {
	var org model.Org
	err := db.Def().Select("id").Where("id = ? AND isdel = 'N'", id).First(&org).Error
	if err != nil {
		return false
	}
	return true
}

// 获取ORG详细信息
func (OrgService) GetOrgByID(id int) (*model.Org, error) {
	var org model.Org
	err := db.Def().Where("id = ? AND isdel = 'N'", id).First(&org).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &org, nil
}
