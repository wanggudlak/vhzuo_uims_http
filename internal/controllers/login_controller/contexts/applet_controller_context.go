package contexts

type AppletCodeLoginResp struct {
	OpenId          string `json:"open_id"`
	IsRegistered    bool   `json:"is_registered"`
	SessionKey      string `json:"session_key"`
	Code            string `json:"code" comment:"UIMS 登录 code"`
	State           string `json:"state"`
	AppletCodeToken string `json:"applet_code_token"`
}

type AppletInfoReq struct {
	AppletCodeToken string `json:"applet_code_token" form:"applet_code_token" binding:"required"`
	Nickname        string `json:"nickname" comment:"昵称"`
	// 性别 M 男, F 女
	Sex      string `json:"sex" comment:"性别"`
	Country  string `json:"country" comment:"国家"`
	Avatar   string `json:"avatar" comment:"头像"`
	Province string `json:"province" comment:"省"`
	City     string `json:"city" comment:"市"`
	// (加密字符串,需要解密)
	EncryptedPhone   string `json:"encrypted_phone" comment:"微信认证的手机号" binding:"required"`
	EncryptedPhoneIv string `json:"encrypted_phone_iv" comment:"手机号解密iv" binding:"required"`
}

type AppletInfoResp struct {
	Code  string `json:"code"`
	State string `json:"state"`
}
