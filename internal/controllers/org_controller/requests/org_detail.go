package requests

type OrgDetailRequest struct {
	ID int `json:"id" form:"id" binding:"required" comment:"组织ID"`
}
