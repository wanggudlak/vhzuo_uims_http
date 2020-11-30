package resource_group_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	requests2 "uims/internal/controllers/resource_group_controller/requests"
	responses2 "uims/internal/controllers/responses"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/db"
)

func List(c *gin.Context) {
	//参数校验处理
	var err error
	var request requests2.ResourceGroupListRequest
	if err = c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}
	// 进行分页数据处理和转换
	page := request.Page
	pagesize := request.PageSize
	if page == 0 {
		page = 1
	}
	if pagesize == 0 {
		pagesize = 10
	}
	page = page - 1

	var source_group_list []model.ResourceGroup
	var count int
	// 根据角色id查询资源组列表
	query := db.Def().Table("uims_res_group").Select("uims_res_group.*")
	if request.RoleID != 0 {
		// 求总数
		counter := query.
			Joins("join uims_role_res_map on uims_role_res_map.res_grp_id = uims_res_group.id").
			Where("uims_role_rs_map.role_id = ? "+
				"and uims_role_res_map.isdel = ? "+
				"and uims_role_res_map.forget_at > ?", request.RoleID, "N", time.Now()).
			Count(&count)
		if counter.Error != nil {
			responses2.Error(c, counter.Error)
			return
		}
		// 求出列表数据
		data := query.
			Where("uims_role_res_map.role_id = ? "+
				"and uims_res_group.client_id = ?", request.RoleID, request.ClientID).
			Offset(page * pagesize).Limit(pagesize).Order("id desc").Scan(&source_group_list)

		if data.Error != nil {
			responses2.Error(c, data.Error)
			return
		}
	} else {
		query = query.Where("isdel = ?", "N").Where("client_id = ?", request.ClientID)

		// 求出总数
		counter := query.Count(&count)
		if counter.Error != nil {
			responses2.Error(c, counter.Error)
			return
		}

		// 求出列表数据
		data := query.Offset(page * pagesize).Limit(pagesize).Order("id desc").Scan(&source_group_list)
		if data.Error != nil {
			responses2.Error(c, data.Error)
			return
		}
	}

	body := map[string]interface{}{
		"data":  source_group_list,
		"total": count,
	}

	responses2.Success(c, "success", body)
}

func Create(c *gin.Context) {
	//进行模型类参数绑定
	var request requests2.ResourceGroupCreateRequest
	var err error
	if err = c.ShouldBindJSON(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	group, err := service.GetResGroupService().CreateResGroup(&request)

	if group == nil && err != nil {
		responses2.Error(c, err)
		return
	}

	responses2.Success(c, "success", group)
}

func Delete(c *gin.Context) {
	var request requests2.ResourceGroupDeleteRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	// delete resource group data
	group, err := service.GetResGroupService().DeleteResGroup(&request)

	if group == nil && err != nil {
		responses2.Error(c, err)
		return
	}

	responses2.Success(c, "success", "nil")
}

func Update(c *gin.Context) {
	// 创建模型绑定参数对象
	var request requests2.ResourceGroupUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	group, err := service.GetResGroupService().UpdateResGroup(&request)

	if group == nil && err != nil {
		responses2.Error(c, err)
		return
	}

	responses2.Success(c, "success", group)
}

func Get(c *gin.Context) {
	var err error
	var request requests2.ResourceGroupIdRequest
	if err = c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}
	var group model.ResourceGroup
	err = db.Def().Where("id = ?", request.ID).First(&group).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}
	fmt.Printf("%v: %T", group)

	var resource_list []model.Resource
	resource_ids := group.ResOfCurr.ResourceIDs
	err = db.Def().Where("id in (?) and isdel = ?", resource_ids, "N").Find(&resource_list).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}
	var resource_data_list []map[string]interface{}
	for _, v := range resource_list {
		data := make(map[string]interface{})
		data["id"] = v.ID
		data["client_id"] = v.ClientId
		data["org_id"] = v.OrgId
		data["res_code"] = v.ResCode
		data["res_front_code"] = v.ResFrontCode
		data["res_type"] = v.ResType
		data["res_name_en"] = v.ResNameEn
		data["res_name_cn"] = v.ResNameCn
		data["res_endp_route"] = v.ResEndpRoute
		data["res_data_location"] = v.ResDataLocation
		data["created_at"] = v.CreatedAt
		resource_data_list = append(resource_data_list, data)
	}
	body := map[string]interface{}{
		"resource_list":       resource_data_list,
		"resource_group_info": group,
	}

	responses2.Success(c, "success", body)
}
