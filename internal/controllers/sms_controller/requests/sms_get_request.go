package requests

type SMSCodeByPhoneRequest struct {
	Scene string `json:"scene" form:"scene" binding:"required" faker:"scene" example:"phone" comment:"场景，用什么媒介发送短信"`
	SPM   string `json:"spm_full_code" form:"spm_full_code" binding:"required" faker:"spm_full_code" example:"fsdafsafsd" comment:"SPM编码"`
	Phone string `json:"phone" form:"phone" binding:"required,len=11,mobile" faker:"phone" example:"13641337591" comment:"手机号"`
}

type SMSCodeByEmailRequest struct {
	Scene string `json:"scene" form:"scene" binding:"required" faker:"scene" example:"email" comment:"场景，用什么媒介发送短信"`
	SPM   string `json:"spm_full_code" form:"spm_full_code" binding:"required" faker:"spm_full_code" example:"fsdafsafsd" comment:"SPM编码"`
	Email string `json:"email" form:"email" binding:"required,email" faker:"email" example:"wanglei@vzhuo.com" comment:"收件人邮箱"`
}

type VerifyRequest struct {
	By   string `json:"by" form:"by" binding:"required" faker:"by" example:"" comment:"手机号或者邮箱"`
	Code string `json:"code" form:"code" binding:"required" faker:"code" example:"code" comment:"验证码"`
	Key  string `json:"key" form:"key" binding:"required" faker:"key" example:"dfadfdfdsfasd" comment:"验证码标识"`
}

type SMSCodeByAccountFormatRequest struct {
	Account      string `json:"account" form:"account" faker:"account" binding:"required,max=32" comment:"账号"`
	Scene        string `json:"scene" form:"scene" binding:"required" faker:"scene" example:"email" comment:"场景，用什么媒介发送短信"`
	SPM          string `json:"spm_full_code" form:"spm_full_code" binding:"required" faker:"spm_full_code" example:"fsdafsafsd" comment:"SPM编码"`
	LoginAccount string `json:"login_account" form:"login_account" faker:"login_account" binding:"required,max=32" comment:"账号"`
}
