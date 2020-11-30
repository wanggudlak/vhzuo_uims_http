package user_role_controller

import (
	"github.com/gin-gonic/gin"
	responses2 "uims/internal/controllers/responses"
	requests2 "uims/internal/controllers/user_role_controller/requests"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/db"
)

// @Summary 为用户添加角色
// @Produce  json
// @Param client_id body int false "客户端ID"
// @Param user_id body int false "用户ID"
// @Param role_id body int true "角色ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/users/role [post]
func AddRole(c *gin.Context) {
	var request requests2.UserRoleNewRequest
	var err error
	if err := c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	maps := make(map[string]int)
	maps["user_id"] = request.UserID
	maps["role_id"] = request.RoleID

	//根据角色id获取client_id
	var role model.Role
	err = db.Def().Where("id = ?", request.RoleID).First(&role).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}
	maps["client_id"] = role.ClientID
	// 查询是否存在 用户关联的角色
	if service.GetUserRoleService().ExistUserRole(maps) {
		responses2.Failed(c, "user role exist", nil)
		return
	}

	// 添加用户角色关联关系
	err = service.GetUserRoleService().AddUserRole(maps)
	if err != nil {
		responses2.Failed(c, err.Error(), nil)
		return
	}

	responses2.Success(c, "success", nil)
}

// @Summary 为用户删除角色
// @Produce  json
// @Param client_id body int false "客户端ID"
// @Param user_id body int false "用户ID"
// @Param role_id body int true "角色ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/users/role [DELETE]
func DeleteRole(c *gin.Context) {
	var request requests2.UserRoleDeleteRequest
	//var err error
	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	maps := make(map[string]int)
	maps["user_id"] = request.UserID
	maps["role_id"] = request.RoleID

	// 查询是否存在 用户关联的角色
	if !service.GetUserRoleService().ExistUserRole(maps) {
		responses2.Failed(c, "user role not exist", nil)
		return
	}

	//根据角色id获取client_id
	var role model.Role
	err1 := db.Def().Where("id = ?", request.RoleID).First(&role).Error
	if err1 != nil {
		responses2.Error(c, err1)
		return
	}
	maps["client_id"] = role.ClientID

	// 删除用户角色关联关系
	if service.GetUserRoleService().DeleteUserRole(maps) != nil {
		responses2.Failed(c, "delete user role failed", nil)
		return
	}

	responses2.Success(c, "success", nil)
}
