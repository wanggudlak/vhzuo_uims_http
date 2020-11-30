package contexts

type BackgroundLoginRequest struct {
	Account string `json:"account" form:"account" binding:"required,max=16" faker:"cc_number" example:"账号" comment:"账号"`
	//用hash加密
	Passwd string `json:"passwd" form:"passwd" binding:"required,max=100,min=6" faker:"password" example:"密码" comment:"密码"`

	//尝试用rsa加解密
	LoginData string `json:"login_data,omitempty" form:"login_data" comment:"登陆数据"`
}
