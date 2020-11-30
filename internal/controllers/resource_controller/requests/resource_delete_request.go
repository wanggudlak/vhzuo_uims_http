package requests

type ResourceDeleteRequest struct {
	ID int `json:"id" form:"id" binding:"required"`
}
