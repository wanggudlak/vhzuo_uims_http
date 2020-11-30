package requests

type UserListRequest struct {
	Page     int    ` form:"page"  binding:"-" json:"page"`
	PageSize int    `form:"pagesize"  binding:"-" json:"page_size"`
	Phone    string `form:"phone" binding:"omitempty" json:"phone"`
	Email    string `form:"email" binding:"omitempty" json:"email"`
}
