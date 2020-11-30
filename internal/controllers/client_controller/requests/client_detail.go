package requests

type ClientDetailRequest struct {
	ID int `json:"id" form:"id" binding:"required" comment:"客户端ID"`
}
