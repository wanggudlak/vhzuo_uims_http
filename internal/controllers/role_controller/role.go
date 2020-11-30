package role_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"uims/conf"
	responses2 "uims/internal/controllers/responses"
	requests2 "uims/internal/controllers/role_controller/requests"
	"uims/internal/service"
	"uims/pkg/tool"
)

// @Summary 获取角色列表
// @Produce  query
// @Param user_id query int true "用户ID"
// @Param org_id query int false "组织ID"
// @Param client_id query int true "客户端ID"
// @Param page query int false "当前页数,默认1"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/roles/list [get]
func List(c *gin.Context) {
	var request requests2.RoleListRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	pageNum := tool.GetPage(c)
	pageSize := conf.AppSetting.PageSize
	maps := make(map[string]interface{})

	if request.UserID != 0 {
		//获取用户的所有角色列表
		userRoleMaps := make(map[string]interface{})
		userRoleMaps["user_id"] = request.UserID

		//查询用户关联的角色, 返回所有角色的id
		RoleIDS, err := service.GetUserRoleService().GetUserRoles(userRoleMaps)
		if err != nil {
			responses2.Failed(c, "get user roles  fail", err)
			return
		}
		maps["role_ids"] = RoleIDS
		maps["isdel"] = "N" // 未软删除

	} else if request.ClientID != 0 {
		// 查询client下的角色列表
		maps["client_id"] = request.ClientID
		maps["isdel"] = "N" // 未软删除

	}
	result, err := service.GetRoleService().GetRoles(pageNum, pageSize, maps)

	if err != nil {
		fmt.Println(err)
		responses2.Failed(c, "get roles list  fail", err)
		return
	}

	total, e := service.GetRoleService().GetRoleTotal(maps)
	if e != nil {
		responses2.Failed(c, "get roles total fail", e)
		return
	}

	data := make(map[string]interface{})
	data["data"] = result
	data["total"] = total

	responses2.Success(c, "success", data)
}

// @Summary 获取角色详情
// @Produce  json
// @Param id query int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/roles [get]
func Detail(c *gin.Context) {
	var request requests2.RoleDetailRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	//查询角色ID 是否存在
	if service.GetRoleService().ExistRoleByID(request.ID) {
		responses2.Failed(c, "id  not exist", nil)
		return
	}

	// 获取角色详细信息
	role, e := service.GetRoleService().GetRole(request.ID)
	if e != nil {
		responses2.Failed(c, "get role id  failed", nil)
		return
	}
	responses2.Success(c, "success", role)
}

// @Summary 添加角色
// @Produce  json
// @Param client_id query int false "ClientID"
// @Param org_id query int false "OrgID"
// @Param NameCN query int true "NameCN"
// @Param NameEN query int false "NameEN"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/roles [post]
func AddRole(c *gin.Context) {
	var request requests2.RoleNewRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}
	roleID := 0
	if service.GetRoleService().ExistRoleByName(request.NameCN, request.ClientID, roleID) {
		responses2.Failed(c, "name_cn exist", nil)
		return
	}
	if service.GetRoleService().AddRole(request.NameCN, request.NameEN, request.ClientID, request.OrgID) != nil {
		responses2.Failed(c, "add role fail", nil)
		return
	}

	responses2.Success(c, "success", nil)
}

// @Summary 更新角色信息
// @Produce  json
// @Param id query int true "ID"
// @Param client_id query int false "ClientID"
// @Param org_id query int false "OrgID"
// @Param NameCN query string false "NameCN"
// @Param NameEN query string false "NameEN"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/roles [put]
func UpdateRole(c *gin.Context) {
	var request requests2.RoleUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	if service.GetRoleService().ExistRoleByID(request.ID) {
		responses2.Failed(c, "id  not exist", nil)
		return
	}

	maps := make(map[string]interface{})
	if request.ClientID > 0 {
		maps["client_id"] = request.ClientID
	}
	if request.OrgID > 0 {
		maps["org_id"] = request.OrgID
	}
	if request.NameCN != "" {
		// 角色中文名字不可重复
		if service.GetRoleService().ExistRoleByName(request.NameCN, request.ClientID, request.ID) {
			responses2.Failed(c, "name_cn exist", nil)
			return
		}
		maps["role_name_cn"] = request.NameCN
	}
	if request.NameEN != "" {
		// 角色英文名字不可重复
		if service.GetRoleService().ExistRoleByNameEN(request.NameEN, request.ClientID, request.ID) {
			responses2.Failed(c, "name_en exist", nil)
			return
		}
		maps["role_name_en"] = request.NameEN
	}

	maps["updated_at"] = time.Now().Format("2006-01-02 15:04:05")

	if service.GetRoleService().UpdateRole(request.ID, maps) != nil {
		responses2.Failed(c, "add role fail", nil)
		return
	}

	responses2.Success(c, "success", nil)
}

// @Summary 删除角色信息 软删除
// @Produce  json
// @Param id query int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/roles [delete]
func DeleteRole(c *gin.Context) {
	var request requests2.RoleDetailRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}
	if service.GetRoleService().ExistRoleByID(request.ID) {
		responses2.Failed(c, "id  not exist", nil)
		return
	}

	if service.GetRoleService().DeleteRole(request.ID) != nil {
		responses2.Failed(c, "delete role fail", nil)
		return
	}

	responses2.Success(c, "success", nil)
}
