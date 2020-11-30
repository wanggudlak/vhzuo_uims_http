package auth_controller

type FindPasswordTokenRequest struct {
	Phone       string `json:"phone" form:"phone" binding:"mobile" comment:"手机号"`
	Email       string `json:"email" form:"email" binding:"" faker:"email" example:"wanglei@vzhuo.com" comment:"收件人邮箱"`
	SMSCode     string `json:"sms_code" form:"sms_code" binding:"required,min=4,max=25" comment:"短信验证码"`
	SPMFullCode string `json:"spm_full_code" form:"spm_full_code" binding:"required,max=255,min=4" comment:"场景码"`
}

type FindPasswordRequest struct {
	FindPasswordToken string `json:"find_password_token" form:"find_password_token" comment:"token"`
	NewPassword       string `json:"new_password" form:"new_password" comment:"新密码"`
}

type RegisterRequest struct {
	Phone       string `json:"phone" form:"phone" comment:"手机号"`
	SMSCode     string `json:"sms_code" form:"sms_code" comment:"短信验证码"`
	Password    string `json:"password" form:"password" comment:"密码"`
	SPMFullCode string `json:"spm_full_code" form:"spm_full_code" binding:"required,max=255,min=4" comment:"场景码"`
}

type RegisterAndLoginRequest struct {
	ClientId    int    `json:"client_id" form:"client_id" binding:"required"`
	RedirectUrl string `json:"redirect_url" form:"redirect_url" binding:"required,url" comment:"重定向地址"`
	State       string `json:"state" form:"state" binding:"max=128" desc:"重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节"`
	Phone       string `json:"phone" form:"phone" comment:"手机号"`
	SMSCode     string `json:"sms_code" form:"sms_code" comment:"短信验证码"`
	Password    string `json:"password" form:"password" comment:"密码"`
	SPMFullCode string `json:"spm_full_code" form:"spm_full_code" binding:"required,max=255,min=4" comment:"场景码"`
}

type GetTokenForFindPasswdFormRequest struct {
	ClientId int `json:"client_id" form:"client_id" binding:"required"`
}
