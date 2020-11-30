package contexts

// LoginHTMLrenderRequest 客户端网站触发登录动作后需要携带的 query params
type LoginHTMLrenderRequest struct {
	SPM         string `form:"spm"`
	RedirectURL string `form:"redirect_url" binding:"required,url" comment:"登录成功后的重定向地址"`
	State       string `form:"state" binding:"max=128" comment:"重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节"`
}
