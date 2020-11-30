package requests

type UserStatusRequest struct {
	Id   int `form:"id" binding:"required"`
	Type int `form:"type" binding:"required"`
}
