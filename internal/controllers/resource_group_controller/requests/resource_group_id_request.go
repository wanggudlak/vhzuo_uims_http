package requests

type ResourceGroupIdRequest struct {
	ID int `form:"id" binding:"required"`
}
