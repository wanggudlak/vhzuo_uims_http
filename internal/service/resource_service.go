package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"uims/internal/controllers/resource_controller/requests"
	"uims/internal/model"
	"uims/pkg/tool"

	"uims/pkg/db"
)

type ResourceService struct{}

type ResourceUpdate struct {
	ResourceOld interface{} `json:"resource_old"`
	ResourceNew interface{} `json:"resource_new"`
}

// ExistResource check existence
func (ResourceService) ExistResource(whereMaps interface{}) bool {
	var resource model.Resource

	err := db.Def().
		Select("id").
		Where("isdel = ?", "N").
		Where(whereMaps).
		First(&resource).
		Error

	if err == nil {
		return true
	}
	return false
}

// check existence not id
func (ResourceService) ExistResourceNotID(id uint, whereMaps interface{}) bool {
	var resource model.Resource

	err := db.Def().
		Select("id").
		Where("id <> ?", id).
		Where("isdel = ?", "N").
		Where(whereMaps).
		First(&resource).
		Error

	if err == nil {
		return true
	}
	return false
}

// GetResourceMapByIDs 查询资源点map
func (ResourceService) GetResourceMapByIDs(ids []int, tx *gorm.DB) ([]model.Resource, error) {
	var resources []model.Resource

	err := tx.
		Select("id, client_id, res_name_en").
		Where("isdel = ?", "N").
		Where("id IN (?)", ids).
		Find(&resources).
		Error

	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		return nil, errors.New("resources not found")
	}

	return resources, nil
}

// Create Resource Data
func (ResourceService) CreateResource(request *requests.ResourceCreateRequest) (*model.Resource, error) {
	// resource data check
	whereMap := map[string]interface{}{
		"client_id":   request.ClientId,
		"res_name_en": request.ResNameEn,
	}
	if GetResourceService().ExistResource(whereMap) {
		return nil, errors.New("resource `" + request.ResNameEn + "` already exists.")
	}

	//生成ResCode和ResFrontCode
	resCode := tool.GenXid()
	resFrontCode := tool.GenXid()

	// 创建数据
	resource := model.Resource{
		ClientId:        request.ClientId,
		OrgId:           request.OrgId,
		ResCode:         resCode,
		ResFrontCode:    resFrontCode,
		ResType:         request.ResType,
		ResSubType:      request.ResSubType,
		ResNameEn:       request.ResNameEn,
		ResNameCn:       request.ResNameCn,
		ResEndpRoute:    request.ResEndpRoute,
		ResDataLocation: &request.ResDataLocation,
		IsDel:           "N",
	}

	// start transaction
	tx := db.Def().Begin()

	err := tx.Create(&resource).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// thrift RPC
	resp := GetThriftClientServer().
		ClientInvoke(int(resource.ClientId), "create_resource", resource)
	if !resp.OK() {
		tx.Rollback()
		return nil, errors.New(resp.Err())
	}

	tx.Commit()

	return &resource, nil
}

// Update Resource Data
func (ResourceService) UpdateResource(request *requests.ResourceUpdateRequest) (*model.Resource, error) {
	var resource model.Resource

	// start transaction
	tx := db.Def().Begin()

	err := tx.Where("id = ?", request.ResourceId).First(&resource).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// resource data check
	whereMap := map[string]interface{}{
		"client_id":   resource.ClientId,
		"res_name_en": request.ResNameEn,
	}
	if GetResourceService().ExistResourceNotID(resource.ID, whereMap) {
		tx.Rollback()
		return nil, errors.New("resource `" + request.ResNameEn + "` already exists.")
	}

	resourceOld := resource

	resource.ResType = request.ResType
	resource.ResSubType = request.ResSubType
	resource.ResNameCn = request.ResNameCn
	resource.ResNameEn = request.ResNameEn
	resource.ResEndpRoute = request.ResEndpRoute
	resource.ResDataLocation = &request.ResDataLocation

	err = tx.Save(&resource).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//thrift RPC
	resourceUpdate := &ResourceUpdate{
		ResourceOld: resourceOld,
		ResourceNew: resource,
	}

	resp := GetThriftClientServer().
		ClientInvoke(int(resource.ClientId), "update_resource", resourceUpdate)
	if !resp.OK() {
		tx.Rollback()
		return nil, errors.New(resp.Err())
	}

	tx.Commit()

	return &resource, nil
}

// Delete Resource Data
func (ResourceService) DeleteResource(request *requests.ResourceDeleteRequest) (*model.Resource, error) {
	var resource model.Resource

	// start transaction
	tx := db.Def().Begin()

	err := tx.Where("id = ? AND isdel = ?", request.ID, "N").First(&resource).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Model(&resource).Where("id = ?", request.ID).Update("isdel", "Y").Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//update resource group data and res_of_curr
	resGroup, err := GetResGroupService().
		DelResNeedUpdateResGroupByResOfCurr(resource.ClientId, request.ID, tx)

	if !resGroup && err != nil {
		tx.Rollback()
		return nil, err
	}

	resp := GetThriftClientServer().
		ClientInvoke(int(resource.ClientId), "delete_resource", resource)
	if !resp.OK() {
		tx.Rollback()
		return nil, errors.New(resp.Err())
	}

	tx.Commit()

	return &resource, nil
}
