package resource_controller

import (
	"github.com/gin-gonic/gin"
	requests2 "uims/internal/controllers/resource_controller/requests"
	responses2 "uims/internal/controllers/responses"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/db"
	"uims/pkg/tool"
)

func Create(c *gin.Context) {
	// 创建模型绑定参数对象
	var request requests2.ResourceCreateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	// create resource data
	resource, err := service.GetResourceService().CreateResource(&request)

	if resource == nil && err != nil {
		responses2.Error(c, err)
		return
	}

	responses2.Success(c, "success", resource)
}

func List(c *gin.Context) {
	// 创建模型绑定参数对象
	var request requests2.ResourceListRequest
	var err error
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
	var count, resource_list interface{}
	var resource_data []model.Resource
	body := make(map[string]interface{})

	if request.UserId != 0 {
		err, resource_list, count = GetResourceByUserOrRole("user", request.UserId, request.ClientId, page, pagesize)
		if err != nil {
			responses2.Error(c, err)
			return
		}
	} else if request.RoleId != 0 {
		// 根据role_id进行资源点的获取
		err, resource_list, count = GetResourceByUserOrRole("user", request.UserId, request.ClientId, page, pagesize)
		if err != nil {
			responses2.Error(c, err)
			return
		}

	} else {
		counter := db.Def().
			Where("isdel = ? and client_id = ?", "N", request.ClientId).
			Find(&resource_data).Count(&count)
		if counter.Error != nil {
			responses2.Error(c, err)
			return
		}
		// 获取分页后的数据
		query := db.Def().Where("isdel = ? and client_id = ?", "N", request.ClientId).
			Offset(page * pagesize).Limit(pagesize).
			Order("id desc").
			Find(&resource_data)
		if query.Error != nil {
			responses2.Error(c, err)
			return
		}
		resource_list = resource_data
	}

	// 组装返回值map
	body = map[string]interface{}{
		"total": count,
		"data":  resource_list,
	}
	responses2.Success(c, "success", body)
}

func Delete(c *gin.Context) {
	var request requests2.ResourceDeleteRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	// delete resource data
	resource, err := service.GetResourceService().DeleteResource(&request)

	if resource == nil && err != nil {
		responses2.Error(c, err)
		return
	}

	responses2.Success(c, "success", "nil")
}

func Get(c *gin.Context) {
	var err error
	var request requests2.ResourceDeleteRequest
	if err = c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	var resource model.Resource

	err = db.Def().Where("id = ?", request.ID).First(&resource).Error

	if err != nil {
		responses2.Error(c, err)
		return
	}
	responses2.Success(c, "success", resource)
}

func Update(c *gin.Context) {
	// 创建模型绑定参数对象
	var request requests2.ResourceUpdateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	// update resource data
	resource, err := service.GetResourceService().UpdateResource(&request)

	if resource == nil && err != nil {
		responses2.Error(c, err)
		return
	}

	responses2.Success(c, "success", resource)
}

func GetResourceByUserOrRole(query_type string, query_id int, client_id int, page int,
	pagesize int) (error, []model.Resource, int) {

	// 定义资源组模型类数据切片
	var res_group_list []model.ResourceGroup
	var err error
	var resource_list []model.Resource
	var count int

	if query_type == "user" {
		// 根据user_id进行资源点的获取
		err = db.Def().Table("uims_res_group").
			Select("uims_res_group.res_of_curr").
			Joins("join uims_role_res_map on uims_role_res_map.res_grp_id = uims_res_group.id").
			Joins("join uims_role on uims_role_res_map.role_id = uims_role.id").
			Joins("join uims_user_role on uims_user_role.role_id = uims_role.id").
			Where("uims_user_role.user_id = ? "+
				"and uims_role.isdel = ? "+
				"and uims_role_res_map.isdel = ? "+
				"and uims_res_group.isdel = ? and uims_res_group.client_id = ?", query_id, "N", "N", "N", client_id).
			Scan(&res_group_list).Error
	} else if query_type == "role" {
		// 根据role_id进行资源点的获取
		err = db.Def().Table("uims_res_group").
			Select("uims_res_group.res_of_curr").
			Joins("join uims_role_res_map on uims_role_res_map.res_grp_id = uims_res_group.id").
			Joins("join uims_role on uims_role_res_map.role_id = uims_role.id").
			Where("uims_role.id = ?"+
				"and uims_role.isdel = ? "+
				"and uims_role_res_map.isdel = ? "+
				"and uims_res_group.isdel = ?"+
				"and uims_res_group.client_id = ?", query_id, "N", "N", "N", client_id).
			Scan(&res_group_list).Error
	}
	if err != nil {
		return err, resource_list, count

	}

	// 定义切片
	var resource_id_list []int
	for _, value := range res_group_list {
		// 所有的切片融合成一个
		resource_id_list = append(resource_id_list, value.ResOfCurr.ResourceIDs...)
	}
	// 进行去重
	resource_id_list = tool.RemoveRepByMap(resource_id_list)

	// 进行in查询获取资源点数据
	counter := db.Def().Where("isdel = ? and id in (?) ", "N", resource_id_list).
		Find(&resource_list).Count(&count)
	if counter.Error != nil {
		return err, resource_list, count
	}
	query := db.Def().Where("isdel = ? and id in (?) ", "N", resource_id_list).
		Offset(page * pagesize).Limit(pagesize).Order("id desc").Find(&resource_list)
	if query.Error != nil {
		return err, resource_list, count
	}

	return err, resource_list, count

}
