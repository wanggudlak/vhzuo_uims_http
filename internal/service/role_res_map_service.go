package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/tool"
)

// 角色关联资源组
type RoleResMapService struct {
}

// 查询角色资源组关联关系表
func (RoleResMapService) ExistRoleResMapByID(id int) bool {
	var roleResMap model.RoleResMap
	err := db.Def().Select("id").Where("id = ? AND isdel = 'N'", id).First(&roleResMap).Error
	if err != nil {
		return false
	}
	return true
}

// 查询角色资源组关联关系表
func (RoleResMapService) ExistRoleResMapByRoleAndGroup(role_id, group_id int) (*model.RoleResMap, error) {
	var roleResMap model.RoleResMap
	err := db.Def().Select("id").Where("role_id = ? AND res_grp_id = ? AND isdel = 'N'", role_id, group_id).First(&roleResMap).Error
	if err != nil && err != gorm.ErrRecordNotFound {

		fmt.Println(err)
		return nil, err
	}
	return &roleResMap, nil
}

// 查询角色关联的所有资源组
func (RoleResMapService) RoleResMapByRoleID(roleID int) ([]model.RoleResMap, error) {
	var roleResMap []model.RoleResMap
	err := db.Def().Where("role_id = ? AND isdel = 'N'", roleID).Order("id desc").Find(&roleResMap).Error
	if err != nil {
		return nil, err
	}
	return roleResMap, nil
}

// 查询资源组关联角色map
func (RoleResMapService) GetRolesMapByGroupId(groupID uint, clientId uint, tx *gorm.DB) ([]model.Role, error) {
	var roles []model.Role

	err := tx.Table("uims_role").
		Select("uims_role.id, uims_role.client_id, uims_role.role_name_en, uims_role.role_name_cn").
		Joins("JOIN uims_role_res_map ON uims_role_res_map.role_id = uims_role.id "+
			"AND uims_role_res_map.isdel = ? "+
			"AND uims_role_res_map.res_grp_id = ? "+
			"AND uims_role_res_map.client_id = ?", "N", groupID, clientId).
		Where("uims_role.isdel = ?", "N").
		Where("uims_role.client_id = ?", clientId).
		Scan(&roles).
		Error

	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, errors.New("roles not found")
	}

	return roles, nil
}

// 查询角色关联资源组map 是否存在
func (RoleResMapService) ExistRoleResMap(id, groupID int) bool {
	var roleResMap model.RoleResMap
	err := db.Def().Select("id").Where("id = ? AND  res_grp_id = ? AND isdel = 'N' AND forget_at > ?",
		id, groupID, time.Now()).First(&roleResMap).Error
	if err != nil {
		return false
	}
	return true
}

// 创建角色与资源组关联关系
func (RoleResMapService) AddRoleResMap(id int, groupID int, forget string) error {
	role, _ := GetRoleService().GetRole(id)

	roleResMap := model.RoleResMap{
		ClientID: role.ClientID,
		ForgetAt: tool.StrToTime(forget),
		ResGrpID: groupID,
		RoleID:   role.ID,
	}
	err := db.Def().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&roleResMap).Error; err != nil {
			return err
		}
		var ids []int
		ids = append(ids, groupID)
		err := RequestUpdateClientRoleRes(tx, id, ids, "add_client_role_res")
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

// 删除角色与资源组关联关系
func (RoleResMapService) DeleteRoleResMap(id int) error {
	err := db.Def().Transaction(func(tx *gorm.DB) error {
		var roleRes model.RoleResMap
		if err := tx.Model(&model.RoleResMap{}).Where("id = ?", id).First(&roleRes).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.RoleResMap{}).Where("id = ?", id).UpdateColumn("isdel", "Y").Error; err != nil {
			return err
		}
		var ids []int
		ids = append(ids, roleRes.ResGrpID)
		err := RequestUpdateClientRoleRes(tx, roleRes.RoleID, ids, "delete_client_role_res")
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

//请求更新子业务系统用户角色资源组数据
func RequestUpdateClientRoleRes(tx *gorm.DB, roleID int, resGroupIDs []int, method string) error {

	type clientRole struct {
		RoleNameCN        string   `json:"role_name_cn"`
		RoleNameEN        string   `json:"role_name_en"`
		ClientFlagCode    string   `json:"client_flag_code"`
		ChangePermissions []string `json:"change_permissions"`
		OnlyItem          bool     `json:"only_item"`
	}
	type resourceIDs struct {
		ResourceIDs []int `json:"resource_ids"`
	}
	role, _ := GetRoleService().GetRole(roleID)

	clientInfo, err := ClientService{}.GetClientByID(role.ClientID)
	if clientInfo == nil || err != nil {
		return err
	}
	//目前只有结算系统有需要回写数据
	//if clientInfo.ClientType != "CASS" {
	//	return nil
	//}
	var resOfCurr, changePermissions []string
	var resourceSlice []int
	tx.Table("uims_res_group").Where("id in (?)", resGroupIDs).Pluck("res_of_curr", &resOfCurr)
	for _, v := range resOfCurr {
		var data resourceIDs
		str := []byte(v)
		_ = json.Unmarshal(str, &data)
		resourceSlice = append(resourceSlice, data.ResourceIDs...)
	}
	resourceSlice = tool.RemoveRepByMap(resourceSlice)
	tx.Table("uims_access_resource").Where("id in (?)", resourceSlice).Pluck("res_name_en", &changePermissions)
	requestRoleResInfo := clientRole{
		RoleNameCN:        role.RoleNameCN,
		RoleNameEN:        role.RoleNameEN,
		ClientFlagCode:    clientInfo.ClientFlagCode,
		ChangePermissions: changePermissions,
		OnlyItem:          true,
	}
	resp := GetThriftClientServer().
		ClientInvoke(role.ClientID, method, requestRoleResInfo)
	if !resp.OK() {
		return errors.New(resp.Err())
	}
	return nil
}

// 删除资源组时，角色与资源组映射关系包含此资源组，需要进行删除
func (RoleResMapService) DelResGroupNeedUpdateRoleResMapByGroupID(clientId uint, groupID uint, tx *gorm.DB) error {
	err := tx.Table("uims_role_res_map").
		Where("isdel = ?", "N").
		Where("client_id = ?", clientId).
		Where("res_grp_id = ?", groupID).
		Update("isdel", "Y").
		Error

	if err != nil {
		return err
	}

	return nil
}

// 删除角色时，角色与资源组映射关系包含此角色，需要进行删除
func (RoleResMapService) DelRoleNeedUpdateRoleResMapByRoleID(clientId int, roleID int, tx *gorm.DB) error {
	err := tx.Table("uims_role_res_map").
		Where("isdel = ?", "N").
		Where("client_id = ?", clientId).
		Where("role_id = ?", roleID).
		Update("isdel", "Y").
		Error

	if err != nil {
		return err
	}

	return nil
}
