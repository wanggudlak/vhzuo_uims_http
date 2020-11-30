package requests

type RoleUpdateRequest struct {
	ID       int    `json:"id" form:"id" binding:"required"`
	ClientID int    `json:"client_id" form:"client_id" binding:"-"`
	OrgID    int    `json:"org_id" form:"org_id" binding:"-"`
	NameEN   string `json:"role_name_en" form:"role_name_en" binding:"-"`
	NameCN   string `json:"role_name_cn" form:"role_name_cn" binding:"-"`
}
