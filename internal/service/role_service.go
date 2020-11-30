package service

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/tool"
)

type RoleService struct {
}

// ExistRoleByID checks if an article exists based on ID
func (RoleService) ExistRoleByID(id int) bool {
	var role model.Role
	err := db.Def().Select("id").Where("id = ? AND isdel = ?", id, "N").First(&role).Error
	if err != nil {
		return true
	}
	return false
}

func (RoleService) GetRole(id int) (*model.Role, error) {
	var role model.Role
	err := db.Def().Where("id = ? AND isdel = 'N'", id).First(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &role, nil
}

// GetRoles gets a list of tags based on paging and constraints
func (RoleService) GetRoles(pageNum int, pageSize int, maps map[string]interface{}) ([]model.Role, error) {
	var (
		roles []model.Role
		err   error
	)

	if _, ok := maps["role_ids"]; ok {
		err = db.Def().Where("id in (?)", maps["role_ids"]).Order("id desc").Find(&roles).Offset(pageNum).Limit(pageSize).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		return roles, nil
	} else {
		err = db.Def().Where(maps).Offset(pageNum).Limit(pageSize).Order("id desc").Find(&roles).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		return roles, nil
	}
}

// GetRoleTotal gets a list of tags based on paging and constraints
func (RoleService) GetRoleTotal(maps map[string]interface{}) (int, error) {
	var count int
	if _, ok := maps["role_ids"]; ok {
		if err := db.Def().
			Where("id in (?)", maps["role_ids"]).
			Model(&model.Role{}).Count(&count).Error; err != nil {
			return 0, err
		}
	} else {
		err := db.Def().Where(maps).Model(&model.Role{}).Count(&count).Error
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

// ExistRoleByName checks if there is a tag with the same name
func (RoleService) ExistRoleByName(name string, clientID int, roleID int) bool {
	var role model.Role
	err := db.Def().Select("id").Where("role_name_cn = ? and isdel = ? and client_id = ? and id != ?", name, "N", clientID, roleID).First(&role).Error
	if err != nil {
		return false
	}
	return true
}

// ExistRoleByNameEn checks if there is a tag with the same name
func (RoleService) ExistRoleByNameEN(name string, clientID int, roleID int) bool {
	var role model.Role
	err := db.Def().Select("id").Where("role_name_en = ? AND isdel = ? and client_id = ? and id != ?", name, "N", clientID, roleID).First(&role).Error
	if err != nil {
		return false
	}
	return true
}

// AddRole Add a Role
func (RoleService) AddRole(nameCn string, nameEn string, clientID int, orgID int) error {
	role := model.Role{
		RoleNameCN: nameCn,
		RoleNameEN: nameEn,
		ClientID:   clientID,
		OrgID:      orgID,
		RoleCode:   "UIMS.SUPERADMIN.001", // 默认值
	}
	var oldRoleInfo model.Role
	err := db.Def().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}

		var deleteResGroupIDs []int
		err := requestUpdateClientRole(tx, role.ID, "add_client_role", deleteResGroupIDs, oldRoleInfo)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// UpdateRole modify a single role
func (RoleService) UpdateRole(id int, data interface{}) error {

	var oldRoleInfo model.Role
	err := db.Def().Where("id = ?", id).First(&oldRoleInfo).Error
	if err != nil {
		return err
	}
	err = db.Def().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Role{}).Where("id = ?", id).Updates(data).Error; err != nil {
			return err
		}

		var deleteResGroupIDs []int
		err := requestUpdateClientRole(tx, id, "update_client_role", deleteResGroupIDs, oldRoleInfo)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// DeleteRole delete a role
func (RoleService) DeleteRole(id int) error {
	var role model.Role

	err := db.Def().Where("id = ? and isdel = 'N'", id).First(&role).Error
	if err != nil {
		return err
	}

	var oldRoleInfo model.Role
	err = db.Def().Transaction(func(tx *gorm.DB) error {
		// 删除角色
		role.IsDel = "Y"
		err := tx.Save(&role).Error
		if err != nil {
			return err
		}

		deleteResGroupIDs := []int{}
		if err = tx.
			Model(&model.RoleResMap{}).
			Where("role_id = ?", id).
			Where("isdel = ?", "N").
			Pluck("res_grp_id", &deleteResGroupIDs).
			Error; err != nil {

			return err
		}

		// 删除角色的关联关系
		err = GetRoleResMapService().
			DelRoleNeedUpdateRoleResMapByRoleID(role.ClientID, role.ID, tx)
		if err != nil {
			return err
		}

		err = requestUpdateClientRole(tx, id, "delete_client_role", deleteResGroupIDs, oldRoleInfo)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

//请求更新子业务系统用户角色数据
func requestUpdateClientRole(tx *gorm.DB, roleID int, method string, deleteResGroupIDs []int, oldRoleInfo model.Role) error {

	type clientRole struct {
		OldRoleNameCN     string   `json:"old_role_name_cn"`
		OldRoleNameEN     string   `json:"old_role_name_en"`
		RoleNameCN        string   `json:"role_name_cn"`
		RoleNameEN        string   `json:"role_name_en"`
		ClientFlagCode    string   `json:"client_flag_code"`
		DeletePermissions []string `json:"delete_permissions"`
		OnlyItem          bool     `json:"only_item"`
	}

	type resourceIDs struct {
		ResourceIDs []int `json:"resource_ids"`
	}
	var role model.Role
	err := tx.Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return err
	}

	clientInfo, err := ClientService{}.GetClientByID(role.ClientID)
	if err != nil {
		return err
	}
	//目前只有结算系统有需要回写数据
	//if clientInfo.ClientType != "CASS" {
	//	return nil
	//}
	var resOfCurr, deletePermissions []string
	var resourceSlice []int
	if deleteResGroupIDs != nil {
		tx.Table("uims_res_group").Where("id in (?)", deleteResGroupIDs).Pluck("res_of_curr", &resOfCurr)
		for _, v := range resOfCurr {
			var data resourceIDs
			str := []byte(v)
			_ = json.Unmarshal(str, &data)
			resourceSlice = append(resourceSlice, data.ResourceIDs...)
		}
		resourceSlice = tool.RemoveRepByMap(resourceSlice)
		tx.Table("uims_access_resource").Where("id in (?)", resourceSlice).Pluck("res_name_en", &deletePermissions)
	}

	requestRoleInfo := clientRole{
		OldRoleNameCN:     oldRoleInfo.RoleNameCN,
		OldRoleNameEN:     oldRoleInfo.RoleNameEN,
		RoleNameCN:        role.RoleNameCN,
		RoleNameEN:        role.RoleNameEN,
		ClientFlagCode:    clientInfo.ClientFlagCode,
		DeletePermissions: deletePermissions,
		OnlyItem:          true,
	}

	resp := GetThriftClientServer().
		ClientInvoke(role.ClientID, method, requestRoleInfo)
	if !resp.OK() {
		return errors.New(resp.Err())
	}
	return nil
}
