package wechatcontroller

import "uims/internal/service/wechatqrcode"

// 微信关注事件
type FollowEventReq struct {
	SceneId int                     `json:"scene_id" binding:"required"`
	WXInfo  wechatqrcode.WeChatInfo `json:"wx_info" binding:"required"`
}
