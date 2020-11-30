package contexts

import (
	"errors"
	"net/url"
	"uims/pkg/encryption"
)

// 登录场景
const (
	AccountPasswdVerifyCodeAuth        = 1 // 用账号、密码、图片验证码请求登录
	AccountPasswdVerifyCodeSMSCodeAuth = 2 // 用账号、密码、图片验证码、手机验证码请求登录
	PhonePasswdAuth                    = 3 // 用手机号、密码请求登录
	EmailPasswdAuth                    = 4 // 用邮箱、密码请求登录
	PhoneVerifyCodeAndSlideCodeAuth    = 5 // 用手机号、滑动验证码+短信验证码请求登录
)

// 请求登录鉴权的基本参数
type BaseAuthReq struct {
	ClientId    int    `json:"client_id" form:"client_id" binding:"required"`
	AuthScene   int    `json:"auth_scene" form:"auth_scene" binding:"required"`
	SPMFullCode string `json:"spm_full_code" form:"spm_full_code" binding:"required,max=255,min=4"`
	RedirectUrl string `json:"redirect_url" form:"redirect_url" binding:"required,url" comment:"重定向地址"`
	State       string `json:"state" form:"state" binding:"max=128" desc:"重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节"`
}

// 用账号、密码、图片验证码请求登录 AccountPasswdVerifyCodeAuth = 1
type AccountPasswdVerifyCodeAuthRequest struct {
	BaseAuthReq
	Account         string `json:"account" form:"account" binding:"required,max=100,min=4" comment:"账号"`
	Password        string `json:"password" form:"password" binding:"required,max=100,min=4" comment:"密码"`
	Code            string `json:"code" form:"code" binding:"required" comment:"验证码"`
	VerificationKey string `json:"verification-key" form:"verification-key" binding:"required,min=10,max=100" comment:"验证码key"`
}

// 用账号、密码、图片验证码、手机验证码请求登录 AccountPasswdVerifyCodeSMSCodeAuth = 2
// 获取手机短信验证码之前需要先校验通过账号密码图片验证码
type AccountPasswdVerifyCodeSMSCodeAuthRequest struct {
	BaseAuthReq
	Account         string `json:"account" form:"account" binding:"omitempty,max=100,min=4" comment:"账号"`
	Password        string `json:"password" form:"password" binding:"omitempty,max=100,min=4" comment:"密码"`
	Code            string `json:"code" form:"code" binding:"omitempty" comment:"验证码"`
	VerificationKey string `json:"verification-key" form:"verification-key" binding:"required,min=10,max=100" comment:"验证码key"`
	PasswordToken   string `json:"password_token" form:"password_token" binding:"omitempty,min=4" comment:"登录一阶段token"`
	SMSCode         string `json:"sms_code" form:"sms_code" binding:"omitempty" comment:"短信验证码"`
}

// 用手机号、密码请求登录 PhonePasswdAuth = 3
type PhonePasswdAuthRequest struct {
	BaseAuthReq
	Phone         string `json:"phone" form:"phone" binding:"required,mobile" comment:"手机号"`
	Password      string `json:"password" form:"password" binding:"required,max=512,min=4" comment:"密码"`
	IsEncryPasswd string `json:"is_encry_passwd" form:"is_encry_passwd" comment:"密码字段的值是否加密了"`
}

// 用邮箱、密码请求登录 EmailPasswdAuth = 4
type EmailPasswdAuthRequest struct {
	BaseAuthReq
	Email         string `json:"email" form:"email" binding:"required" comment:"邮箱"`
	Password      string `json:"password" form:"password" binding:"required,max=100,min=4" comment:"密码"`
	IsEncryPasswd string `json:"is_encry_passwd" form:"is_encry_passwd" comment:"密码字段的值是否加密了"`
}

// 用手机号、短信验证码请求登录 PhoneVerifyCodeAuth = 5
// 获取短信验证码之前通过滑动验证码校验通过之后才发送短信
type PhoneVerifyCodeAndSlideCodeAuthRequest struct {
	BaseAuthReq
	Phone   string `json:"phone" form:"phone" binding:"required,mobile" comment:"手机号"`
	SMSCode string `json:"sms_code" form:"sms_code" comment:"短信验证码"`
	//Scene 			string `json:"scene" form:"scene" binding:"required" faker:"scene" example:"email" comment:"场景"`
	VerificationKey string `json:"verification-key" form:"verification-key" binding:"required,min=10,max=100" comment:"验证码key"`
}

type HandlerReqAfterValidate interface {
	decryPasswdIfEncrypted() error
}

func (req *PhonePasswdAuthRequest) DecryPasswdIfEncrypted() error {
	if "yes" == req.IsEncryPasswd {
		decryBytes, err := encryption.DecryptWithRSA([]byte(req.Password))
		if err != nil {
			return err
		}
		req.Password = string(decryBytes)
	}
	return nil
}

type AuthenticateResp struct {
	RedirectUrl string `json:"redirect_url"`
	AuthCode    string `json:"code"`
}

func ParseAndPrepareAuthResponse(redirectURL, authCode, state string) (*AuthenticateResp, error) {
	var err error
	pURLParsed := &url.URL{}
	if pURLParsed, err = url.Parse(redirectURL); err != nil {
		return nil, errors.New("回调地址格式解析失败: " + err.Error())
	}

	URLvalues := pURLParsed.Query()
	URLvalues.Add("state", state)
	URLvalues.Add("code", authCode)
	pURLParsed.RawQuery = URLvalues.Encode()

	return &AuthenticateResp{
		AuthCode:    authCode,
		RedirectUrl: pURLParsed.String(),
	}, nil
}

type WeChatQRCodeReq struct {
	BaseAuthReq
}

type WeChatQRCodeLoginReq struct {
	BaseAuthReq
	SceneId int `json:"scene_id" form:"scene_id" comment:"场景id"`
}

type WeChatQRCodeLoginResp struct {
	RedirectURL    string `json:"redirect_url" comment:"登录成功跳转地址"`
	Timeout        bool   `json:"timeout" comment:"是否超时"`
	AuthOK         bool   `json:"auth_ok" comment:"是否验证成功"`
	NeedBindPhone  bool   `json:"need_bind_phone" comment:"是否需要绑定手机号"`
	BindPhoneToken string `json:"bind_phone_token" comment:"绑定手机号 token"`
}

type BindPhoneReq struct {
	BaseAuthReq
	BindPhoneToken string `json:"bind_phone_token" binding:"required" comment:"token"`
	Phone          string `json:"phone" form:"phone" binding:"required,mobile" comment:"手机号"`
	SMSCode        string `json:"sms_code" form:"sms_code" comment:"短信验证码"`
}
