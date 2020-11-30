package login_controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"uims/conf"
	"uims/internal/controllers/login_controller/contexts"
	resp "uims/internal/controllers/responses"
	"uims/internal/controllers/sms_controller"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/jwtauth"
	"uims/internal/service/uuid"
	"uims/internal/service/wechat"
	"uims/internal/service/wechatqrcode"
	"uims/pkg/db"
	"uims/pkg/gjwt"
	"uims/pkg/glog"
	"uims/pkg/randc"
	thriftclient "uims/pkg/thrift/client"
	"uims/pkg/tool"
	"uims/pkg/wechatserver"
)

// 场景
// 用户在扫码后, 页面将获取到 微信的 code + state
// js 提交 code + state , 该接口将进行登录操作
// 登录完毕后, 将返回需要的 redirect_url 带 code 参数地址
func WeChatCodeLogin(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	data := c.Query("data")
	glog.Channel("casswechat").WithFields(log.Fields{
		"code":  code,
		"state": state,
		"data":  data,
	}).Info("用户扫码后带来code,status信息,将换取微信信息")
	var redirectUrl = ""
	err := func() error {
		b, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return fmt.Errorf("解析 data (%s) 参数失败: %s", data, err.Error())
		}
		type params struct {
			SMPFullCode string `json:"smp_full_code"`
			RedirectUrl string `json:"redirect_url"`
			ClientId    int    `json:"client_id"`
		}
		var p params
		err = json.Unmarshal(b, &p)
		if err != nil {
			return fmt.Errorf("un json data (%s) 参数失败: %s", string(b), err.Error())
		}
		config, err := wechat.GetConfig(p.ClientId)
		if err != nil {
			return err
		}
		token, err := wechatserver.Cli(*config).ExchangeToken(code)
		if err != nil {
			return err
		}
		userInfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
		if err != nil {
			return err
		}
		glog.Channel("casswechat").WithFields(log.Fields{
			"code":      code,
			"state":     state,
			"user_info": userInfo,
		}).Info("解析到的用户微信信息")
		// 找 openId 对应用户
		userWeChat := &model.UserWeChat{}
		err = db.Def().
			Where(&model.UserWeChat{
				WeChatOpenId: userInfo.OpenId,
			}).
			First(&userWeChat).Error
		glog.Channel("casswechat").WithFields(log.Fields{
			"code":        code,
			"state":       state,
			"open_id":     userInfo.OpenId,
			"user_wechat": userWeChat,
		}).Info("openId 对应用户信息")
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return fmt.Errorf("该微信未绑定用户")
			}
			return err
		}
		if userWeChat.WeChatUnionId == "" && userInfo.UnionId != "" {
			// 更新 union_id
			userWeChat.WeChatUnionId = userInfo.UnionId
			err = db.Def().Save(&userWeChat).Error
			if err != nil {
				return err
			}
		}
		user := &model.User{}
		err = db.Def().Where("id = ?", userWeChat.UserId).First(&user).Error
		glog.Channel("casswechat").WithFields(log.Fields{
			"code":    code,
			"state":   state,
			"open_id": userInfo.OpenId,
			"user":    user,
		}).Info("openId 对应用户信息")
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return fmt.Errorf("该微信未绑定用户")
			}
			return err
		}
		// 进行登录操作
		userJwtService := jwtauth.UserJwtAuth{
			OpenId:   user.OpenID,
			ClientId: uint(p.ClientId),
			Account:  user.Account,
			State:    state,
		}
		if userJwtService.IsFreeze() {
			return jwtauth.FreezeErr
		}
		code := userJwtService.GenerateCode()
		if u, err := url.Parse(p.RedirectUrl); err != nil {
			return err
		} else {
			q := u.Query()
			q.Add("state", state)
			q.Add("code", code)
			u.RawQuery = q.Encode()
			redirectUrl = u.String()
		}
		return nil
	}()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, redirectUrl)
	return
}

// 获取关注二维码
func WeChatQRCode(c *gin.Context) {
	var err error
	var req contexts.WeChatQRCodeReq
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	type P struct {
		SceneId   int    `json:"scene_id"`
		TicketUrl string `json:"ticket_url"`
	}
	p := P{}
	r := service.Response{}
	service.ThriftClientServer{}.InvokeMP(service.Request{
		BRequest: thriftclient.BRequest{
			MethodName: "get_uims_wx_qr",
			Params:     map[string]string{},
		},
	}, &r)
	if !r.OK() {
		resp.Failed(c, r.Err(), nil)
	}
	err = r.ParseContent(&p)
	if err != nil {
		resp.Error(c, err)
		return
	}
	err = wechatqrcode.SetScene(p.SceneId)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Success(c, "", gin.H{
		"qr_code":  p.TicketUrl,
		"scene_id": p.SceneId,
	})
	return
}

// 获取扫码结果, 需要轮询调用
// 参数  scene_id
//
// 响应
//
// redirect_url 登录成功跳转地址
// timeout 二维码是否超时， 超时后需要刷新
// auth_ok 是否登录成功, 用来判断能否跳转了
// need_bind_phone 是否需要绑定手机号
// bind_phone_token 绑定手机号接口需要使用的token
// @Summary 获取扫码登录结果, 需要轮询调用
// @Produce  query
// @Param data body contexts.WeChatQRCodeLoginReq
// @Success 200 {object} contexts.WeChatQRCodeLoginResp
// @Router /api/login/wechat/qr_code/login [GET]
func WeChatQRCodeLogin(c *gin.Context) {
	var err error
	var req contexts.WeChatQRCodeLoginReq
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	var user model.User
	var redirectURL string
	err = func() error {
		// 判断 scene 是否存在
		if !wechatqrcode.Exists(req.SceneId) {
			resp.Success(c, "", contexts.WeChatQRCodeLoginResp{
				RedirectURL:    "",
				Timeout:        true,
				AuthOK:         false,
				NeedBindPhone:  false,
				BindPhoneToken: "",
			})
			return nil
		}

		// 判断是否已经扫码
		if !wechatqrcode.IsScan(req.SceneId) {
			resp.Success(c, "", contexts.WeChatQRCodeLoginResp{
				RedirectURL:    "",
				Timeout:        false,
				AuthOK:         false,
				NeedBindPhone:  false,
				BindPhoneToken: "",
			})
			return nil
		}

		// 判断是否需要绑定手机号
		if wechatqrcode.NeedBindPhone(req.SceneId) {
			w, err := wechatqrcode.GetNeedBindPhoneWeChatInfo(req.SceneId)
			if err != nil {
				return errors.Wrap(err, "获取缓存中的微信数据失败")
			}
			claims := wechatqrcode.BindPhoneClaims{
				WeChatInfo: *w,
			}
			token, err := gjwt.CreateToken(&claims)
			if err != nil {
				return err
			}
			resp.Success(c, "", contexts.WeChatQRCodeLoginResp{
				RedirectURL:    "",
				Timeout:        false,
				AuthOK:         false,
				NeedBindPhone:  true,
				BindPhoneToken: token,
			})
			return nil
		}

		// 判断是否存在 user_id
		if userID, err := wechatqrcode.GetSceneUserID(req.SceneId); err != nil {
			return err
		} else {
			err = db.Def().Where("id = ?", userID).First(&user).Error
			if err != nil {
				return err
			}
			// 进行登录操作
			userJwtService := jwtauth.UserJwtAuth{
				OpenId:   user.OpenID,
				ClientId: uint(req.ClientId),
				Account:  user.Account,
				State:    req.State,
			}
			if userJwtService.IsFreeze() {
				return jwtauth.FreezeErr
			}
			code := userJwtService.GenerateCode()
			if u, err := url.Parse(req.RedirectUrl); err != nil {
				return err
			} else {
				q := u.Query()
				q.Add("state", req.State)
				q.Add("code", code)
				u.RawQuery = q.Encode()
				redirectURL = u.String()
				resp.Success(c, "", contexts.WeChatQRCodeLoginResp{
					RedirectURL:    redirectURL,
					Timeout:        false,
					AuthOK:         true,
					NeedBindPhone:  false,
					BindPhoneToken: "",
				})
				return nil
			}
		}
	}()
	if err != nil {
		resp.Error(c, err)
		return
	}
	return
}

// @Summary 微信扫码登录绑定手机号
// @Produce  json
// @Param data body contexts.BindPhoneReq
// @Success 200 {string} json "{"code":200,"data":{"redirect_url": "http://host?cdoe=123&state=abc"},"msg":"ok"}"
// @Router /api/login/wechat/qr_code/bind/phone [POST]
func WeChatQRCodeLoginBindPhone(c *gin.Context) {
	var err error
	var req contexts.BindPhoneReq
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	var redirectURL string
	var user model.User
	var userWeChat model.UserWeChat
	err = func() error {
		// 校验验证码
		scene := sms_controller.UsePhoneBindWeChat
		cacheCode := service.RedisDefaultCache.RedisGetStr(service.MakeSMSRedisCacheKey(req.Phone, req.SPMFullCode+scene))
		if conf.Switch.SMSCaptcha && cacheCode != req.SMSCode {
			return fmt.Errorf("短信验证码错误")
		}
		// 解析 token
		claims := wechatqrcode.BindPhoneClaims{}
		err = gjwt.Parse(req.BindPhoneToken, &claims)
		if err != nil {
			return errors.Wrap(err, "解析 bind phone token 失败")
		}
		config, err := wechat.GetConfig(req.ClientId)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("客户端ID: %d 查询微信配置失败", req.ClientId))
		}
		// 判断手机号是否被使用
		err = db.Def().Where("phone = ?", req.Phone).First(&user).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return errors.Wrap(err, "查询 user 信息失败")
			}
		} else {
			// 手机号已经存在的情况下, 判断是否已经绑定微信
			userWeChat := model.UserWeChat{}
			err = db.Def().Where("user_id = ?", user.ID).
				Where("wechat_id = ?", config.WeChatId).
				First(&userWeChat).Error
			if err != nil {
				if !gorm.IsRecordNotFoundError(err) {
					return errors.Wrap(err, "查询 user_wechat 信息失败")
				}
			} else {
				return errors.New("该手机号已绑定其它微信")
			}
		}
		// 判断微信是否被使用
		err = db.Def().Where("wechat_open_id = ?", claims.WeChatInfo.OpenID).First(&userWeChat).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return errors.Wrap(err, "查询 user_wechat 数据失败")
			}
		} else {
			return errors.New("该微信已绑定账号")
		}
		// 未注册的用户进行注册
		sex := ""
		if claims.WeChatInfo.Sex == 1 {
			sex = "M"
		}
		if claims.WeChatInfo.Sex == 2 {
			sex = "F"
		}
		err = db.Def().Transaction(func(tx *gorm.DB) error {
			if user.ID == 0 {
				user = model.User{
					OpenID:   randc.UUID(),
					UserType: "VDK",
					UserCode: tool.GenerateUuid(),
					NaCode:   "+86",
					Phone:    &req.Phone,
				}
				err = db.Def().Create(&user).Error
				if err != nil {
					return errors.Wrap(err, "保存 user 数据失败")
				}
				var userInfo = model.UserInfo{
					UserCode: uuid.GenerateForUIMS().String(),
					Phone:    *user.Phone,
					NaCode:   user.NaCode,
					Sex:      sex,
					UserID:   user.ID,
					UserType: user.UserType,
				}
				err = db.Def().Create(&userInfo).Error
				if err != nil {
					return errors.Wrap(err, "保存 user_info 失败")
				}
			}
			// 进行绑定
			newUserWeChat := model.UserWeChat{
				UserId:       uint(user.ID),
				WeChatId:     config.WeChatId,
				Nickname:     claims.WeChatInfo.Nickname,
				Sex:          sex,
				Country:      claims.WeChatInfo.Country,
				Avatar:       claims.WeChatInfo.HeadImgURL,
				Privilege:    "",
				Province:     claims.WeChatInfo.Province,
				City:         claims.WeChatInfo.City,
				WeChatOpenId: claims.WeChatInfo.OpenID,
			}
			err = db.Def().Save(&newUserWeChat).Error
			if err != nil {
				return errors.Wrap(err, "保存 user_wechat 数据失败")
			}
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "创建用户失败")
		}

		// 进行登录
		// 进行登录操作
		userJwtService := jwtauth.UserJwtAuth{
			OpenId:   user.OpenID,
			ClientId: uint(req.ClientId),
			Account:  user.Account,
			State:    req.State,
		}
		if userJwtService.IsFreeze() {
			return jwtauth.FreezeErr
		}
		code := userJwtService.GenerateCode()
		if u, err := url.Parse(req.RedirectUrl); err != nil {
			return err
		} else {
			q := u.Query()
			q.Add("state", req.State)
			q.Add("code", code)
			u.RawQuery = q.Encode()
			redirectURL = u.String()
		}
		return nil
	}()
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Success(c, "", gin.H{
		"redirect_url": redirectURL,
	})
	return
}
