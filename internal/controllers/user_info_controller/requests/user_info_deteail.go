package requests

type UserDetailRequest struct {
	Id int `form:"id" binding:"required"`
}
