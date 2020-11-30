package requests

type RoleDetailRequest struct {
	ID int `form:"id" json:"id" binding:"required" commit:"角色ID"`
}
