package requests

type RoleListRequest struct {
	Page     int `form:"page" json:"" binding:"required"`
	ClientID int `form:"client_id" json:"client_id" binding:"-"`
	OrgID    int `form:"org_id" json:"org_id" binding:"-"`
	UserID   int `form:"user_id" json:"user_id" binding:"-"`
}
