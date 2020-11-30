package oauthcontroller

type AccessTokenReq struct {
	Code      string `json:"code" binding:"required,len=32"`
	GrantType string `json:"grant_type" binding:"required,eq=authorization_code"`
}

type AccessTokenResp struct {
	AccessToken  string `json:"access_token" comment:"凭据"`
	ExpiresIn    uint   `json:"expires_in" comment:"有效时间"`
	RefreshToken string `json:"refresh_token" comment:"刷新 access_token 凭据"`
	OpenId       string `json:"open_id" comment:"用户唯一 id"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" comment:"刷新 access_token 凭据"`
}

type UserInfoReq struct {
	AccessToken string `json:"access_token" comment:"凭据"`
	OpenId      string `json:"open_id" comment:"用户唯一 id"`
}

type UserInfoWeChat struct {
	Nickname string `json:"nickname"`
	OpenID   string `json:"open_id"`
	Avatar   string `json:"avatar"`
	UUID     string `json:"uuid"`
}

type UserInfoResp struct {
	OpenId   string           `json:"open_id"`
	Account  string           `json:"account"`
	UserCode string           `json:"user_code"`
	Phone    string           `json:"phone"`
	Email    string           `json:"email"`
	Nickname string           `json:"nickname"`
	WeChats  []UserInfoWeChat `json:"wechats"`
}

type GetWeChatBindURLReq struct {
	AccessToken string `json:"access_token" binding:"required" comment:"用户凭据"`
	WeChatUUID  string `json:"wechat_uuid" binding:"required" comment:"微信主体唯一标识"`
	RedirectURL string `json:"redirect_url" binding:"required" comment:"绑定完成重定向地址"`
	State       string `json:"state" binding:"omitempty,max=128" comment:"业务场景"`
}

type GetWeChatBindURLResp struct {
	URL string `json:"url"`
}

type GetWeChaUnbindURLReq struct {
	AccessToken string `json:"access_token" binding:"required" comment:"用户凭据"`
	WeChatUUID  string `json:"wechat_uuid" binding:"required" comment:"微信主体唯一标识"`
	RedirectURL string `json:"redirect_url" binding:"required" comment:"解除绑定重定向地址"`
	State       string `json:"state" binding:"omitempty,max=128" comment:"业务场景"`
}

type GetWeChatUnbindURLResp struct {
	URL string `json:"url"`
}

type UnbindWeChatReq struct {
	AccessToken string `json:"access_token" binding:"required" comment:"用户凭据"`
	WeChatUUID  string `json:"wechat_uuid" binding:"required" comment:"微信主体唯一标识"`
}

type UnbindWeChatResp struct {
}

type IsBindWeChatReq struct {
	AccessToken string `json:"access_token" binding:"required" comment:"用户凭据"`
	WeChatUUID  string `json:"wechat_uuid" binding:"required" comment:"微信主体唯一标识"`
}

type IsBindWeChatResp struct {
	IsBind bool `json:"is_bind"`
}
