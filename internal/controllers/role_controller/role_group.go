package role_controller

import (
	"github.com/gin-gonic/gin"
	responses2 "uims/internal/controllers/responses"
	requests2 "uims/internal/controllers/role_controller/requests"
	"uims/internal/model"
	"uims/internal/service"
)

// @Summary 获取角色的所有资源组详情
// @Produce  form
// @Param id query int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/roles/group [POST]
func GroupDetail(c *gin.Context) {
	var request requests2.RoleDetailRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	// 查询角色ID
	if service.GetRoleService().ExistRoleByID(request.ID) {
		responses2.Failed(c, "id  not exist", nil)
		return
	}

	// 查询角色关联的所有资源组
	roleResMaps, err := service.GetRoleResMapService().RoleResMapByRoleID(request.ID)
	if err != nil {
		responses2.Failed(c, "select res group relationship exist ", nil)
		return
	}

	body := make(map[string]interface{})
	// 组合数据
	var result []model.ResourceGroup
	for _, roleResMap := range roleResMaps {
		resGroup, err := service.GetResGroupService().GetResGroupByID(roleResMap.ResGrpID)

		if err != nil {
			responses2.Failed(c, "res group not exist ", nil)
			return
		}
		result = append(result, *resGroup)
	}
	body["data"] = result
	body["total"] = len(result)

	responses2.Success(c, "success", body)
}

// @Summary 角色更新资源组
// @Produce  json
// @Param id body int true "ID"
// @Param group_id body int true "GroupID"
// @Param forget query datetime true "Forget"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/roles/group [POST]
func AddGroup(c *gin.Context) {
	var request requests2.GroupNewRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	// 查询角色
	if service.GetRoleService().ExistRoleByID(request.ID) {
		responses2.Failed(c, "role id  not exist", nil)
		return
	}

	// 查询资源组
	//if service.GetResGroupService().ExistResGroupByID(request.ID) {
	//	responses2.Failed(c, "res group not exist", nil)
	//	return
	//}

	// 查询角色资源组 是否关联 已经关联的 不能再次关联
	if service.GetRoleResMapService().ExistRoleResMap(request.ID, request.GroupID) {
		responses2.Failed(c, "res group relationship exist", nil)
		return
	}
	// TODO 已经关联的 过期了的,修改过期时间  可优化

	if request.Forget == "" {
		request.Forget = "3000-01-01"
	}

	// 关联资源组
	if service.GetRoleResMapService().AddRoleResMap(request.ID, request.GroupID, request.Forget) != nil {
		responses2.Failed(c, "add role fail", nil)
		return
	}

	responses2.Success(c, "success", nil)
}

// @Summary 删除角色与资源组的关系
// @Produce  json
// @Param role_id body int true "role_id"
// @Param group_id body int true "group_id"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/roles/group [DELETE]
func DeleteGroup(c *gin.Context) {
	var request requests2.RoleGroupMapRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	// 查询角色资源组关联关系
	roleResMaps, err := service.GetRoleResMapService().ExistRoleResMapByRoleAndGroup(request.RoleID, request.GroupID)
	if err != nil {
		responses2.Failed(c, "res group relationship  not exist", nil)
		return
	}

	// 删除角色资源组关联关系  软删除
	if service.GetRoleResMapService().DeleteRoleResMap(roleResMaps.ID) != nil {
		responses2.Failed(c, "res group relationship exist", nil)
		return
	}

	responses2.Success(c, "success", nil)
}
