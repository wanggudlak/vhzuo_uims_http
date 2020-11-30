package requests

type UserRoleNewRequest struct {
	UserID int `form:"user_id" json:"user_id" binding:"required" commit:"用户ID"`
	RoleID int `form:"role_id" json:"role_id" binding:"required" commit:"角色ID"`
}
