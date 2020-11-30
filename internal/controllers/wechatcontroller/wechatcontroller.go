package wechatcontroller

import (
	"encoding/json"
	"fmt"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	resp "uims/internal/controllers/responses"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/jwtauth"
	"uims/internal/service/wechat"
	"uims/pkg/db"
	"uims/pkg/gjwt"
	"uims/pkg/tool"
	"uims/pkg/wechatserver"
)

// 进行绑定微信操作
func Bind(c *gin.Context) {
	var err error
	var req BindReq
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	var user model.User
	var cache = BindCache{}
	var u *url.URL
	err = func() error {
		cacheStr := service.RedisDefaultCache.RedisGetStr(service.RedisKey(req.State))
		if cacheStr == "" {
			return fmt.Errorf("无效的 state 值")
		}

		err = json.Unmarshal([]byte(cacheStr), &cache)
		if err != nil {
			return err
		}

		u, err = url.Parse(cache.RedirectURL)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		resp.Error(c, err)
		return
	}
	q := u.Query()

	err = func() error {
		err = db.Def().Where("id = ?", cache.UserId).First(&user).Error
		if err != nil {
			return err
		}
		userWeChat := model.UserWeChat{}
		err = db.Def().
			Select("id").
			Where(&model.UserWeChat{
				WeChatId: cache.WeChatId,
				UserId:   cache.UserId,
			}).First(&userWeChat).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return err
			}
		} else {
			return fmt.Errorf("该用户已绑定微信")
		}
		config, err := wechat.GetConfigByWeChatId(uint(cache.WeChatId))
		if err != nil {
			return err
		}
		token, err := wechatserver.Cli(*config).ExchangeToken(req.Code)
		if err != nil {
			return err
		}
		userInfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
		if err != nil {
			return err
		}
		err = db.Def().
			Select("id").
			Where("wechat_open_id = ?", userInfo.OpenId).
			Or("wechat_union_id = ?", userInfo.UnionId).
			First(&userWeChat).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return err
			}
		} else {
			return fmt.Errorf("该微信已绑定用户")
		}
		// 进行绑定操作
		sex := ""
		if userInfo.Sex == 1 {
			sex = "F"
		}
		if userInfo.Sex == 2 {
			sex = "M"
		}
		userWeChat = model.UserWeChat{
			UserId:        uint(user.ID),
			WeChatId:      cache.WeChatId,
			Nickname:      userInfo.Nickname,
			Sex:           sex,
			Country:       userInfo.Country,
			Avatar:        userInfo.HeadImageURL,
			Privilege:     tool.JSONString(userInfo.Privilege),
			Province:      userInfo.Province,
			City:          userInfo.City,
			WeChatOpenId:  userInfo.OpenId,
			WeChatUnionId: userInfo.UnionId,
		}
		err = db.Def().Create(&userWeChat).Error
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		// 失败跳转反馈信息
		q.Add("state", cache.State)
		q.Add("msg", err.Error())
		q.Add("code", "")
		u.RawQuery = q.Encode()
		c.Redirect(http.StatusFound, u.String())
		return
	}
	// 进行登录操作
	userJwtService := jwtauth.UserJwtAuth{
		OpenId:   user.OpenID,
		ClientId: cache.ClientId,
		Account:  user.Account,
		State:    cache.State,
	}
	if userJwtService.IsFreeze() {
		resp.Error(c, jwtauth.FreezeErr)
		return
	}
	code := userJwtService.GenerateCode()
	q.Add("state", cache.State)
	q.Add("code", code)
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusFound, u.String())
	return

}

// 进行解绑微信操作
func Unbind(c *gin.Context) {
	var err error
	var req UnbindReq
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	var user model.User
	var u *url.URL
	var claims UnbindClaims
	err = func() error {
		err = gjwt.Parse(req.Token, &claims)
		if err != nil {
			return errors.Wrap(err, "解析token失败")
		}
		u, err = url.Parse(claims.RedirectURL)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		resp.Error(c, err)
		return
	}
	q := u.Query()

	err = func() error {
		err = db.Def().Where("id = ?", claims.UserId).First(&user).Error
		if err != nil {
			return err
		}
		userWeChat := model.UserWeChat{}
		err = db.Def().
			Select("id").
			Where(&model.UserWeChat{
				WeChatId: claims.WeChatId,
				UserId:   claims.UserId,
			}).First(&userWeChat).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return fmt.Errorf("该用户未绑定微信")
			}
			return errors.Wrap(err, "查询 user_wechat 失败")
		}
		err = db.Def().Unscoped().Delete(&userWeChat).Error
		if err != nil {
			return errors.Wrap(err, "解除绑定失败")
		}
		return nil
	}()
	if err != nil {
		// 失败跳转反馈信息
		q.Add("state", claims.State)
		q.Add("msg", err.Error())
		q.Add("code", "")
		u.RawQuery = q.Encode()
		c.Redirect(http.StatusFound, u.String())
		return
	}
	// 进行登录操作
	userJwtService := jwtauth.UserJwtAuth{
		OpenId:   user.OpenID,
		ClientId: claims.ClientId,
		Account:  user.Account,
		State:    claims.State,
	}
	if userJwtService.IsFreeze() {
		resp.Error(c, jwtauth.FreezeErr)
		return
	}
	code := userJwtService.GenerateCode()
	q.Add("state", claims.State)
	q.Add("code", code)
	u.RawQuery = q.Encode()
	c.Redirect(http.StatusFound, u.String())
	return
}
