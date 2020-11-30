package requests

type ResourceListRequest struct {
	Page     int `form:"page"  binding:"-" json:"page"`
	PageSize int `form:"pagesize"  binding:"-" json:"pagesize"`
	RoleId   int `form:"role_id" binding:"-" json:"role_id"`
	UserId   int `form:"user_id" binding:"-" json:"user_id"`
	ClientId int `form:"client_id"  binding:"required" json:"client_id"`
}
