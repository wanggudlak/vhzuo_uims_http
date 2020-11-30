package service

import (
	"github.com/jinzhu/gorm"
	"uims/internal/model"
	"uims/pkg/db"
)

type UserOrgService struct {
}

func (UserOrgService) SaveUserOrgData(userID int, orgName string, roleName string, tx *gorm.DB) error {
	//超级管理员没有组织概念
	if roleName == "super_admin" {
		return nil
	}
	var orgInfo model.Org
	err := db.Def().Model("uims_organization").Where("org_name_cn = ? ", orgName).First(&orgInfo).Error
	if err != nil {
		return err
	}
	var roleInfo model.Role
	err = db.Def().Model("uims_role").Where("role_name_en = ? ", roleName).First(&roleInfo).Error
	if err != nil {
		return err
	}

	var userOrgInfo model.UserOrg
	var userOrg = model.UserOrg{
		UserID:   userID,
		ClientID: roleInfo.ClientID,
		OrgID:    orgInfo.ID,
	}
	if err := tx.Where("user_id = ? and client_id = ?", userID, roleInfo.ClientID).
		First(&userOrgInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := tx.Create(&userOrg).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if err := tx.Model(&userOrg).
			Where("id = ?", userOrgInfo.ID).
			Update(&userOrg).Error; err != nil {
			return err
		}
	}

	return nil
}
