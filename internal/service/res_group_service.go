package service

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"reflect"
	"strconv"
	"uims/internal/controllers/resource_group_controller/requests"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/slices"
	"uims/pkg/tool"
)

type ResGroupService struct{}

type Items struct {
	Roles []RoleWithResource `json:"roles"`
}

type RoleWithResource struct {
	ID           int              `json:"id"`
	ClientID     int              `json:"client_id"`
	RoleNameEN   string           `json:"role_name_en"`
	RoleNameCN   string           `json:"role_name_cn"`
	DelResources []model.Resource `json:"del_resources"`
	AddResources []model.Resource `json:"add_resources"`
}

type RoleRelationResource struct {
	RoleID    int        `gorm:"column:role_id;" json:"role_id"`
	ResGrpID  int        `gorm:"column:res_grp_id;" json:"res_grp_id"`
	ResOfCurr *ResOfCurr `gorm:"column:res_of_curr;" json:"res_of_curr"`
}

type ResOfCurr struct {
	ResourceIDs []int `gorm:"column:resource_ids" json:"resource_ids"`
}

// Scan 实现方法
func (r *ResOfCurr) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &r)
}

func (ResGroupService) ExistResGroupByID(id int) bool {
	var group model.ResourceGroup
	err := db.Def().Select("id").Where("id = ? AND isdel = 'N'", id).First(&group).Error
	if err != nil {
		return false
	}
	return true
}

// check existence
func (ResGroupService) ExistResGroup(whereMaps interface{}) bool {
	var group model.ResourceGroup

	err := db.Def().
		Select("id").
		Where("isdel = ?", "N").
		Where(whereMaps).
		First(&group).
		Error

	if err == nil {
		return true
	}
	return false
}

// check existence not id
func (ResGroupService) ExistResGroupNotID(id uint, whereMaps interface{}) bool {
	var group model.ResourceGroup

	err := db.Def().
		Select("id").
		Where("id <> ?", id).
		Where("isdel = ?", "N").
		Where(whereMaps).
		First(&group).
		Error

	if err == nil {
		return true
	}
	return false
}

func (ResGroupService) GetResGroupByID(id int) (*model.ResourceGroup, error) {
	var resGroup model.ResourceGroup
	err := db.Def().Where("id = ? AND isdel = 'N'", id).First(&resGroup).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &resGroup, nil
}

//  GetResGroupByClientId checks if there is a tag with the same clientId
func (ResGroupService) GetResGroupByClientId(clientId uint) ([]model.ResourceGroup, error) {
	var resGroup []model.ResourceGroup

	err := db.Def().
		Table("uims_res_group").
		Select("id,res_group_en,res_group_cn,res_of_curr,client_id,org_id").
		Where("isdel = ?", "N").
		Where("client_id = ?", clientId).
		Scan(&resGroup).
		Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return resGroup, nil
}

// check resource_ids
func (ResGroupService) CheckResGroupByResourceIds(resourceIds []int, clientId uint) error {

	for _, id := range resourceIds {
		whereMap := map[string]interface{}{
			"id":        id,
			"client_id": clientId,
		}

		if GetResourceService().ExistResource(whereMap) == false {
			return errors.New("id " + strconv.Itoa(id) + " not found")
		}
	}

	return nil
}

// Create Resource Group Data
func (ResGroupService) CreateResGroup(request *requests.ResourceGroupCreateRequest) (*model.ResourceGroup, error) {
	// resource group data check
	whereMap := map[string]interface{}{
		"client_id":    request.ClientId,
		"res_group_en": request.ResGroupEn,
	}
	if GetResGroupService().ExistResGroup(whereMap) {
		return nil, errors.New("resource group `" + request.ResGroupEn + "` already exists.")
	}

	// check resource_ids
	if len(request.ResOfCurr.ResourceIDs) == 0 {
		return nil, errors.New("`res_of_curr` is null")
	} else {
		err := GetResGroupService().CheckResGroupByResourceIds(request.ResOfCurr.ResourceIDs, request.ClientId)
		if err != nil {
			return nil, err
		}
	}

	resGroupCode := tool.GenXid()

	group := model.ResourceGroup{
		ResGroupCode: resGroupCode,
		ResGroupCn:   request.ResGroupCn,
		ResGroupEn:   request.ResGroupEn,
		ResGroupType: request.ResGroupType,
		ResOfCurr:    request.ResOfCurr,
		ClientId:     request.ClientId,
		OrgId:        request.OrgId,
	}

	err := db.Def().Create(&group).Error

	if err != nil {
		return nil, err
	}

	return &group, nil
}

// Update Resource Group Data
func (ResGroupService) UpdateResGroup(request *requests.ResourceGroupUpdateRequest) (*model.ResourceGroup, error) {
	var group model.ResourceGroup

	// start transaction
	tx := db.Def().Begin()

	err := tx.Where("id = ?", request.ID).First(&group).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// resource group data check
	whereMap := map[string]interface{}{
		"client_id":    group.ClientId,
		"res_group_en": request.ResGroupEn,
	}
	if GetResGroupService().ExistResGroupNotID(group.ID, whereMap) {
		tx.Rollback()
		return nil, errors.New("resource group `" + request.ResGroupEn + "` already exists.")
	}

	// check resource_ids
	if len(request.ResOfCurr.ResourceIDs) == 0 {
		tx.Rollback()
		return nil, errors.New("`res_of_curr` is null")
	} else {
		err = GetResGroupService().CheckResGroupByResourceIds(request.ResOfCurr.ResourceIDs, group.ClientId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	groupOld := group

	group.ResGroupEn = request.ResGroupEn
	group.ResGroupCn = request.ResGroupCn
	group.ResGroupType = request.ResGroupType
	group.ResOfCurr = request.ResOfCurr

	resGroup, err := GetResGroupService().UpdateResGroupNeedThriftClient(&groupOld, &group, tx)
	if resGroup == false && err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Save(&group).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &group, nil
}

//更新资源组时，需要thrift RPC
func (ResGroupService) UpdateResGroupNeedThriftClient(groupOld *model.ResourceGroup, group *model.ResourceGroup, tx *gorm.DB) (bool, error) {
	var updateResourceIDs []int
	for _, id := range group.ResOfCurr.ResourceIDs {
		updateResourceIDs = append(updateResourceIDs, id)
	}

	if reflect.DeepEqual(groupOld.ResOfCurr.ResourceIDs, group.ResOfCurr.ResourceIDs) {
		return true, nil
	}

	// roles
	roles, err := GetRoleResMapService().GetRolesMapByGroupId(group.ID, group.ClientId, tx)
	if err != nil {
		return true, nil
	}

	if len(group.ResOfCurr.ResourceIDs) == 0 {
		return true, nil
	}

	roleWithResource := make([]RoleWithResource, len(roles))

	for k, val := range roles {
		delResourceIDs, addResourceIDs, err := GetResGroupService().
			FilterResourceIDs("update", val.ID, int(group.ID), groupOld.ResOfCurr.ResourceIDs, group.ResOfCurr.ResourceIDs, tx)

		if delResourceIDs == nil && addResourceIDs == nil && err != nil {
			continue
		}

		roleWithResource[k].ID = val.ID
		roleWithResource[k].ClientID = val.ClientID
		roleWithResource[k].RoleNameEN = val.RoleNameEN
		roleWithResource[k].RoleNameCN = val.RoleNameCN

		if delResourceIDs != nil {
			delResources, err := GetResourceService().GetResourceMapByIDs(delResourceIDs, tx)
			if delResources != nil && err == nil {
				roleWithResource[k].DelResources = delResources
			} else {
				roleWithResource[k].DelResources = nil
			}
		}

		if addResourceIDs != nil {
			addResources, err := GetResourceService().GetResourceMapByIDs(addResourceIDs, tx)
			if addResources != nil && err == nil {
				roleWithResource[k].AddResources = addResources
			} else {
				roleWithResource[k].AddResources = nil
			}
		}
	}

	if len(roleWithResource) == 0 {
		return true, nil
	}

	// thrift RPC
	items := &Items{
		Roles: roleWithResource,
	}

	// 指针已修改值，此处需要还原值
	group.ResOfCurr.ResourceIDs = updateResourceIDs

	resp := GetThriftClientServer().ClientInvoke(int(group.ClientId), "update_resource_group", items)
	if !resp.OK() {
		return false, errors.New(resp.Err())
	}

	return true, nil
}

//删除资源点时，资源组数据包含此资源点，需要更新资源组res_of_curr字段信息
func (ResGroupService) DelResNeedUpdateResGroupByResOfCurr(clientId uint, resId int, tx *gorm.DB) (bool, error) {
	resGroup, err := GetResGroupService().GetResGroupByClientId(clientId)

	if resGroup != nil && err == nil {
		for _, item := range resGroup {
			if isExist, _ := slices.IsExistValue(resId, item.ResOfCurr.ResourceIDs); isExist {
				resourceIDs := slices.RemoveIntSlice(item.ResOfCurr.ResourceIDs, resId)
				res, err := GetResGroupService().UpdateResGroupByResOfCurr(uint(item.ID), resourceIDs, tx)
				if !res && err != nil {
					return false, err
				}
			}
		}
	}

	return true, nil
}

// update resource group res_of_curr
func (ResGroupService) UpdateResGroupByResOfCurr(id uint, resOfCurr []int, tx *gorm.DB) (bool, error) {
	var group model.ResourceGroup

	err := tx.Where("id = ?", id).First(&group).Error
	if err != nil {
		return false, err
	}

	group.ResOfCurr.ResourceIDs = resOfCurr

	err = tx.Save(&group).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// delete resource group data
func (ResGroupService) DeleteResGroup(request *requests.ResourceGroupDeleteRequest) (*model.ResourceGroup, error) {
	var group model.ResourceGroup

	// start transaction
	tx := db.Def().Begin()

	err := tx.Where("id = ? AND isdel = ?", request.ID, "N").First(&group).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Model(&group).Where("id = ?", request.ID).Update("isdel", "Y").Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	resGroup, err := GetResGroupService().DeleteResGroupNeedThriftClient(&group, tx)

	if resGroup == false && err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &group, nil
}

//删除资源组时，需要thrift RPC
func (ResGroupService) DeleteResGroupNeedThriftClient(group *model.ResourceGroup, tx *gorm.DB) (bool, error) {
	// roles
	roles, err := GetRoleResMapService().GetRolesMapByGroupId(group.ID, group.ClientId, tx)
	if err != nil {
		return true, nil
	}

	if len(group.ResOfCurr.ResourceIDs) == 0 {
		return true, nil
	}

	roleWithResource := make([]RoleWithResource, len(roles))

	for k, val := range roles {
		delResourceIDs, _, err := GetResGroupService().
			FilterResourceIDs("delete", val.ID, int(group.ID), nil, group.ResOfCurr.ResourceIDs, tx)

		// 删除本身自己
		if delResourceIDs == nil && err == nil {
			delResourceIDs = group.ResOfCurr.ResourceIDs
		}

		// 过滤之后，无需删除
		if delResourceIDs == nil && err != nil {
			continue
		}

		roleWithResource[k].ID = val.ID
		roleWithResource[k].ClientID = val.ClientID
		roleWithResource[k].RoleNameEN = val.RoleNameEN
		roleWithResource[k].RoleNameCN = val.RoleNameCN

		resources, err := GetResourceService().GetResourceMapByIDs(delResourceIDs, tx)
		if resources != nil && err == nil {
			roleWithResource[k].DelResources = resources
		} else {
			roleWithResource[k].DelResources = nil
		}
	}

	// delete uims_role_res_map
	err = GetRoleResMapService().DelResGroupNeedUpdateRoleResMapByGroupID(group.ClientId, group.ID, tx)
	if err != nil {
		return false, err
	}

	if len(roleWithResource) == 0 {
		return true, nil
	}

	// thrift RPC
	items := &Items{
		Roles: roleWithResource,
	}

	resp := GetThriftClientServer().ClientInvoke(int(group.ClientId), "delete_resource_group", items)
	if !resp.OK() {
		return false, errors.New(resp.Err())
	}
	return true, nil
}

// 过滤 ResourceIDs
func (ResGroupService) FilterResourceIDs(handle string, roleID int, groupID int, oldResourceIDs []int, resourceIDs []int, tx *gorm.DB) ([]int, []int, error) {
	var (
		roleRelationResource []RoleRelationResource
		addResourceIDs       []int
		delResourceIDs       []int
	)

	querySql := tx.Table("uims_role_res_map").
		Select("uims_role_res_map.role_id, uims_role_res_map.res_grp_id, uims_res_group.res_of_curr").
		Joins("LEFT JOIN uims_res_group ON uims_role_res_map.res_grp_id = uims_res_group.id "+
			"AND uims_res_group.isdel = ?", "N").
		Where("uims_role_res_map.isdel = ?", "N").
		Where("uims_role_res_map.role_id = ?", roleID)

	// delete
	if handle == "delete" {
		querySql = querySql.Where("uims_role_res_map.res_grp_id <> ?", groupID)
	}

	err := querySql.Scan(&roleRelationResource).Error

	if len(roleRelationResource) == 0 || err != nil {
		return nil, nil, err
	}

	for _, value := range roleRelationResource {
		for _, val := range value.ResOfCurr.ResourceIDs {
			if handle == "delete" { // delete
				if isExist, _ := slices.IsExistValue(val, resourceIDs); isExist {
					resourceIDs = slices.RemoveIntSlice(resourceIDs, val)
				}
			} else { // update
				if isExist, _ := slices.IsExistValue(val, resourceIDs); isExist {
					resourceIDs = slices.RemoveIntSlice(resourceIDs, val)
					oldResourceIDs = slices.RemoveIntSlice(oldResourceIDs, val)
				}

				if value.ResGrpID != groupID {
					if isExist, _ := slices.IsExistValue(val, oldResourceIDs); isExist {
						oldResourceIDs = slices.RemoveIntSlice(oldResourceIDs, val)
					}
				}
			}
		}

		if handle == "delete" { // delete
			delResourceIDs = append(delResourceIDs, resourceIDs...)
		} else { // update
			addResourceIDs = append(addResourceIDs, resourceIDs...)
			delResourceIDs = append(delResourceIDs, oldResourceIDs...)
		}
	}

	if handle == "delete" { // delete
		if len(delResourceIDs) == 0 {
			return nil, nil, errors.New("not delete resource")
		}
	} else { // update
		if len(addResourceIDs) == 0 && len(delResourceIDs) == 0 {
			return nil, nil, errors.New("not update resource")
		}
	}

	return tool.RemoveRepByMap(delResourceIDs), tool.RemoveRepByMap(addResourceIDs), nil
}
