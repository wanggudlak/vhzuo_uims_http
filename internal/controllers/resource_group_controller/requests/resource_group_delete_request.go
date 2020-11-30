package requests

type ResourceGroupDeleteRequest struct {
	ID int `json:"id" form:"id" binding:"required"`
}
