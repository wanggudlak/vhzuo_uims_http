package requests

type RoleNewRequest struct {
	ClientID int    `json:"client_id" form:"client_id" binding:"required"`
	OrgID    int    `json:"org_id" form:"org_id" binding:"-"`
	NameEN   string `json:"role_name_en" form:"role_name_en" binding:"required"`
	NameCN   string `json:"role_name_cn" form:"role_name_cn" binding:"required"`
}
