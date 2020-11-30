package requests

type UserStoreRequest struct {
	Email       string `json:"email" form:"email" binding:"required" faker:"email" example:"admin@uims.com" comment:"邮箱"`
	Account     string `json:"account" form:"account" binding:"required,max=16" faker:"cc_number" example:"账号" comment:"账号"`
	Phone       string `json:"phone" form:"phone" binding:"required,len=11,mobile" example:"13517210606" comment:"手机号"`
	Passwd      string `json:"passwd" form:"passwd" binding:"required,max=100,min=6" faker:"password" example:"密码" comment:"密码"`
	EncryptType int    `json:"encrypt_type" form:"encrypt_type" binding:"gte=0"`
}
