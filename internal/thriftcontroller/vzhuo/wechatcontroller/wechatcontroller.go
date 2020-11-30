package wechatcontroller

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"uims/internal/model"
	"uims/internal/service/wechatqrcode"
	"uims/pkg/db"
	thriftserver "uims/pkg/thrift/server"
)

// 微信用户关注事件
func FollowEvent(c *thriftserver.Context) {
	var err error
	var req FollowEventReq
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var userWeChat model.UserWeChat

	err = func() error {
		// 查看 Scene 是否设置
		if !wechatqrcode.Exists(req.SceneId) {
			return errors.New("scene_id 未设置")
		}
		// 设置为已扫码
		err = wechatqrcode.SetScanOK(req.SceneId)
		if err != nil {
			return errors.Wrap(err, "设置 scan_ok 失败")
		}
		// 查 openID 有没有对应用户
		err = db.Def().Where("wechat_open_id = ?", req.WXInfo.OpenID).First(&userWeChat).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return errors.Wrap(err, "查询 user_wechat 数据失败")
			}
			// 未绑定手机号
			err = wechatqrcode.SetNeedBindPhone(req.SceneId, req.WXInfo)
			if err != nil {
				return errors.Wrap(err, "设置 need_bind_phone 失败")
			}
		} else {
			// 已绑定手机号
			err = wechatqrcode.SetSceneAuthOK(req.SceneId, userWeChat.UserId)
			if err != nil {
				return errors.Wrap(err, "设置 auth_ok 失败")
			}
		}
		return nil
	}()
	if err != nil {
		c.Response.Error(err)
		return
	}

	c.Response.Success(nil, "ok")
	return
}
