package oauthcontroller

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	url2 "net/url"
	"time"
	"uims/conf"
	"uims/internal/controllers/wechatcontroller"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/jwtauth"
	wechat2 "uims/internal/service/wechat"
	"uims/pkg/db"
	"uims/pkg/gjwt"
	"uims/pkg/randc"
	thriftserver "uims/pkg/thrift/server"
)

// 通过 code 换取 access_token
// code 使用完后将会删除
// 有效期默认为 2h
func AccessToken(c *thriftserver.Context) {
	var err error
	var req AccessTokenReq
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var jwtAuth jwtauth.UserJwtAuth
	var accessToken string
	var refreshToken string
	err = func() error {
		err = jwtAuth.ParseCode(req.Code)
		jwtAuth.RemoveCode(req.Code)
		if err != nil {
			return err
		}
		accessToken, err = jwtAuth.GenerateAccessToken()
		if err != nil {
			return err
		}
		refreshToken, err = jwtAuth.GenerateRefreshToken()
		return nil
	}()

	if err != nil {
		c.Response.Error(err)
		return
	}

	c.Response.Success(AccessTokenResp{
		AccessToken:  accessToken,
		ExpiresIn:    uint(jwtauth.AccessTokenTTL.Seconds()),
		RefreshToken: refreshToken,
		OpenId:       jwtAuth.OpenId,
	}, "")
}

// 通过 refresh_token 换取 access_token
func RefreshToken(c *thriftserver.Context) {
	var err error
	var req RefreshTokenReq
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var refreshClaims jwtauth.RefreshClaims
	err = gjwt.Parse(req.RefreshToken, &refreshClaims)
	if err != nil {
		c.Response.Error(err)
		return
	}

	auth := jwtauth.UserJwtAuth{
		ClientId: refreshClaims.ClientId,
		OpenId:   refreshClaims.OpenId,
	}

	if auth.IsFreeze() {
		c.Response.Error(jwtauth.FreezeErr)
		return
	}

	accessToken, err := auth.GenerateAccessToken()
	if err != nil {
		c.Response.Error(err)
		return
	}

	c.Response.Success(AccessTokenResp{
		AccessToken:  accessToken,
		ExpiresIn:    uint(jwtauth.AccessTokenTTL.Seconds()),
		RefreshToken: req.RefreshToken,
		OpenId:       auth.OpenId,
	}, "")
}

// access_token 获取用户信息
func UserInfo(c *thriftserver.Context) {
	var err error
	var req UserInfoReq
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var resp UserInfoResp
	err = func() error {
		var claims jwtauth.AccessClaims
		err = gjwt.Parse(req.AccessToken, &claims)
		if err != nil {
			return err
		}
		var u = model.User{}
		err = db.Def().
			Select([]string{"id", "open_id", "account", "user_code", "phone", "email"}).
			Where(&model.User{OpenID: claims.OpenId}).
			First(&u).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return errors.Wrap(err, "未找到用户信息")
			}
			return errors.Wrap(err, "查询 user 失败")
		}

		var uInfo model.UserInfo
		err = db.Def().Select("nickname").Where("user_id = ?", u.ID).First(&uInfo).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return errors.Wrap(err, "查询 user_info 失败")
			}
		}
		resp.OpenId = u.OpenID
		resp.Account = u.Account
		resp.UserCode = u.UserCode
		resp.Phone = *u.Phone
		resp.Email = u.Email
		resp.Nickname = uInfo.Nickname
		var weChats []UserInfoWeChat
		err := db.Def().
			Table("uims_user_wechat as uw").
			Select([]string{"uw.nickname", "uw.wechat_open_id as open_id", "uw.avatar", "w.uuid"}).
			Joins("left join uims_wechat as w on w.id = uw.wechat_id").
			Where("uw.user_id = ?", u.ID).
			Scan(&weChats).Error
		if err != nil {
			return err
		}
		resp.WeChats = weChats
		return nil
	}()
	if err != nil {
		c.Response.Error(err)
		return
	}
	c.Response.Success(resp, "")
	return
}

// 获取绑定微信的链接, 客户端需要跳转到这个连接去进行绑定
func GetBindWeChatURL(c *thriftserver.Context) {
	var err error
	var req GetWeChatBindURLReq
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var url = ""
	err = func() error {
		claims := jwtauth.AccessClaims{}
		err = gjwt.Parse(req.AccessToken, &claims)
		if err != nil {
			return err
		}
		wechat := model.WeChat{}
		err = db.Def().Where("uuid = ?", req.WeChatUUID).First(&wechat).Error
		if err != nil {
			return err
		}
		user := model.User{}
		err = db.Def().
			Select("id").
			Where("open_id = ?", claims.OpenId).
			First(&user).Error
		if err != nil {
			return err
		}
		userWeChat := model.UserWeChat{}
		err = db.Def().
			Select("id").
			Where(&model.UserWeChat{
				WeChatId: wechat.ID,
				UserId:   uint(user.ID),
			}).First(&userWeChat).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return err
			}
		} else {
			return fmt.Errorf("用户已绑定微信")
		}
		config, err := wechat2.GetConfig(int(claims.ClientId))
		if err != nil {
			return err
		}
		redirectURL, err := url2.Parse(conf.URL + "/api/oauth/wechat/bind")
		if err != nil {
			return err
		}
		state := "bindwechat:" + randc.UUID()
		var cache = wechatcontroller.BindCache{
			ClientId:    claims.ClientId,
			WeChatId:    wechat.ID,
			UserId:      uint(user.ID),
			RedirectURL: req.RedirectURL,
			State:       req.State,
		}
		cacheJSON, _ := json.Marshal(&cache)
		err = service.RedisDefaultCache.RedisCacheString(service.RedisKey(state), string(cacheJSON), 5*time.Minute)
		if err != nil {
			return err
		}
		url = fmt.Sprintf("https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
			config.AppId,
			url2.QueryEscape(redirectURL.String()),
			state,
		)
		return nil
	}()
	if err != nil {
		c.Response.Error(err)
		return
	}
	c.Response.Success(GetWeChatBindURLResp{URL: url}, "")
	return
}

// 获取解绑微信地址
func GetUnbindWeChatURL(c *thriftserver.Context) {
	var err error
	var req GetWeChaUnbindURLReq
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var url = ""
	err = func() error {
		claims := jwtauth.AccessClaims{}
		err = gjwt.Parse(req.AccessToken, &claims)
		if err != nil {
			return err
		}
		wechat := model.WeChat{}
		err = db.Def().Where("uuid = ?", req.WeChatUUID).First(&wechat).Error
		if err != nil {
			return errors.Wrap(err, "wechat 查询失败")
		}
		user := model.User{}
		err = db.Def().
			Select("id").
			Where("open_id = ?", claims.OpenId).
			First(&user).Error
		if err != nil {
			return errors.Wrapf(err, "open_id: %s 对应用户不存在", claims.OpenId)
		}
		userWeChat := model.UserWeChat{}
		err = db.Def().
			Select("id").
			Where(&model.UserWeChat{
				WeChatId: wechat.ID,
				UserId:   uint(user.ID),
			}).First(&userWeChat).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return fmt.Errorf("用户未绑定微信")
			}
			return errors.Wrap(err, "查询 user_wechat 失败")
		}
		if err != nil {
			return err
		}
		redirectURL, err := url2.Parse(conf.URL + "/api/oauth/wechat/unbind")
		if err != nil {
			return err
		}
		unbindClaims := wechatcontroller.UnbindClaims{
			ClientId:    claims.ClientId,
			WeChatId:    wechat.ID,
			UserId:      uint(user.ID),
			RedirectURL: req.RedirectURL,
			State:       req.State,
		}
		unbindToken, err := gjwt.CreateToken(&unbindClaims)
		if err != nil {
			return errors.Wrap(err, "创建 token 失败")
		}
		v := redirectURL.Query()
		v.Add("token", unbindToken)
		redirectURL.RawQuery = v.Encode()
		// 跳转到解绑地址
		url = redirectURL.String()
		return nil
	}()
	if err != nil {
		c.Response.Error(err)
		return
	}
	c.Response.Success(GetWeChatUnbindURLResp{URL: url}, "")
	return
}

func IsBindWeChat(c *thriftserver.Context) {
	var err error
	var req IsBindWeChatReq
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var isBind bool
	err = func() error {
		claims := jwtauth.AccessClaims{}
		err = gjwt.Parse(req.AccessToken, &claims)
		if err != nil {
			return err
		}
		wechat := model.WeChat{}
		err = db.Def().Where("uuid = ?", req.WeChatUUID).First(&wechat).Error
		if err != nil {
			return errors.Wrap(err, "wechat 查询失败")
		}
		user := model.User{}
		err = db.Def().
			Select("id").
			Where("open_id = ?", claims.OpenId).
			First(&user).Error
		if err != nil {
			return errors.Wrapf(err, "open_id: %s 对应用户不存在", claims.OpenId)
		}
		userWeChat := model.UserWeChat{}
		err = db.Def().
			Select("id").
			Where(&model.UserWeChat{
				WeChatId: wechat.ID,
				UserId:   uint(user.ID),
			}).First(&userWeChat).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				isBind = false
				return nil
			}
			return errors.Wrap(err, "查询 user_wechat 失败")
		}
		isBind = true
		return nil
	}()
	if err != nil {
		c.Response.Error(err)
		return
	}
	c.Response.Success(IsBindWeChatResp{IsBind: isBind}, "")
	return
}
