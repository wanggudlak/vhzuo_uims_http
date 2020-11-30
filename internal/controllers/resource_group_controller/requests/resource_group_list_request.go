package requests

type ResourceGroupListRequest struct {
	Page     int `form:"page"  binding:"-" json:"page"`
	RoleID   int `form:"role_id" binding:"-" json:"role_id"`
	UserId   int `form:"user_id" binding:"-" json:"user_id"`
	PageSize int `form:"pagesize"  binding:"-" json:"page_size"`
	ClientID int `form:"client_id"  binding:"required" json:"client_id"`
}
