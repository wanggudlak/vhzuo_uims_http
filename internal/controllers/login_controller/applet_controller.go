package login_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/medivhzhan/weapp/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	"uims/internal/controllers/login_controller/contexts"
	resp "uims/internal/controllers/responses"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/jwtauth"
	"uims/internal/service/uuid"
	"uims/internal/service/wechat"
	"uims/pkg/db"
	"uims/pkg/gjwt"
	"uims/pkg/glog"
	"uims/pkg/randc"
)

type AppletCodeLoginClaims struct {
	gjwt.Jwt
	OpenId      string
	SessionKey  string
	State       string
	SPMFullCode string
	ClientId    uint
	WeChatId    uint
	UnionId     string
}

// @Summary 小程序登录第一步, wechat code换取 UIMS code
// @Produce  json
// @Param code query string true "微信小程序授权登录code"
// @Param state query string true "业务码, 业务系统自定义, 将原样返回"
// @Success 200 {object} AppletCodeLoginResp
// @Router /api/login/applet/code [get]
func AppletCodeLogin(c *gin.Context) {
	var err error
	err = func() error {
		code := c.Query("code")
		state := c.Query("state")
		smpFullCode := c.Query("smp_full_code")
		glog.Channel("applet").Printf("code: %s state : %s smp: %s \n", code, state, smpFullCode)
		if code == "" {
			return fmt.Errorf("code 参数必须传入")
		}
		if smpFullCode == "" {
			return fmt.Errorf("smp_full_code 参数必须传入")
		}
		p := service.ParseSPMstring(smpFullCode)
		clientSetting := model.ClientSetting{}
		err = db.Def().Where("spm_full_code = ?", p.FullCode).First(&clientSetting).Error
		if err != nil {
			glog.Channel("applet").Printf("查询 client_setting 失败: %s", err.Error())
			return err
		}
		config, err := wechat.GetConfig(int(clientSetting.ClientID))
		if err != nil {
			return err
		}
		loginRes, err := weapp.Login(config.AppId, config.Secret, code)
		glog.Channel("applet").Printf("小程序code登录从微信获取到信息: %+v", loginRes)
		if err != nil {
			return err
		}
		if err := loginRes.GetResponseError(); err != nil {
			return err
		}
		userWeChat := model.UserWeChat{}
		err = db.Def().
			Select("user_id, wechat_open_id, wechat_union_id").
			Where(&model.UserWeChat{
				WeChatOpenId: loginRes.OpenID,
			}).
			First(&userWeChat).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				glog.Channel("applet").WithFields(log.Fields{
					"request": loginRes,
				}).Error("user_wechat 信息未查询到")
				return errors.Wrap(err, "user_wechat 信息未查询到")
			}
			// 未查到记录, 则需要下一步获取微信授权信息再登录
			j := AppletCodeLoginClaims{
				UnionId:     loginRes.UnionID,
				OpenId:      loginRes.OpenID,
				SessionKey:  loginRes.SessionKey,
				State:       state,
				SPMFullCode: smpFullCode,
				ClientId:    config.ClientId,
				WeChatId:    config.WeChatId,
			}
			j.SetIssue()
			j.SetAudience("applet_login")
			j.SetTTL(5 * time.Minute)
			token, err := gjwt.CreateToken(&j)
			if err != nil {
				return err
			}
			resp.Success(c, "", contexts.AppletCodeLoginResp{
				OpenId:          loginRes.OpenID,
				IsRegistered:    false,
				SessionKey:      loginRes.SessionKey,
				State:           state,
				AppletCodeToken: token,
			})
			return nil
		}
		if userWeChat.WeChatUnionId == "" && loginRes.UnionID != "" {
			// 更新 union_id
			userWeChat.WeChatUnionId = loginRes.UnionID
			err = db.Def().Save(&userWeChat).Error
			if err != nil {
				return err
			}
		}
		user := model.User{}
		err = db.Def().Where("id = ?", userWeChat.UserId).First(&user).Error
		if err != nil {
			glog.Channel("applet").WithFields(log.Fields{
				"user_wechat": userWeChat,
			}).Errorf("查询 user %d 失败: %s", userWeChat.UserId, err.Error())
			return errors.Wrap(err, "查询 user 失败")
		}
		// 生成授权码
		userJwtService := jwtauth.UserJwtAuth{
			OpenId:   user.OpenID,
			ClientId: clientSetting.ClientID,
			Account:  user.Account,
			State:    state,
		}
		if userJwtService.IsFreeze() {
			return jwtauth.FreezeErr
		}
		authCode := userJwtService.GenerateCode()
		resp.Success(c, "", contexts.AppletCodeLoginResp{
			OpenId:          loginRes.OpenID,
			IsRegistered:    true,
			SessionKey:      loginRes.SessionKey,
			State:           state,
			AppletCodeToken: "",
			Code:            authCode,
		})
		return nil
	}()
	if err != nil {
		resp.Error(c, err)
		return
	}
}

// @Summary 小程序登录第二步, 提交授权信息进行注册
// @Produce  json
// @Param user body AppletInfoReq
// @Success 200 {object} AppletInfoResp
// @Router /api/login/applet/info [post]
func AppletInfo(c *gin.Context) {
	var err error
	var req contexts.AppletInfoReq
	var authCode string
	var state string
	var phone string
	var naCode string
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	err = func() error {
		claims := AppletCodeLoginClaims{}
		err = gjwt.Parse(req.AppletCodeToken, &claims)
		if err != nil {
			return err
		}
		// 解密手机号
		decryptMobileRes, err := weapp.DecryptMobile(claims.SessionKey, req.EncryptedPhone, req.EncryptedPhoneIv)
		if err != nil {
			return err
		}
		phone = decryptMobileRes.PhoneNumber
		naCode = "+" + decryptMobileRes.CountryCode
		state = claims.State
		// 检查用户是否存在
		userWeChat := model.UserWeChat{}
		err = db.Def().
			Where(&model.UserWeChat{
				WeChatOpenId:  claims.OpenId,
				WeChatUnionId: claims.UnionId,
			}).
			First(&userWeChat).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return err
			}
		} else {
			return fmt.Errorf("用户已经注册, 请重新发起登录授权")
		}
		// 进行注册
		var user model.User
		err = db.Def().Transaction(func(tx *gorm.DB) error {
			p := service.ParseSPMstring(claims.SPMFullCode)
			var client model.Client
			err = db.Def().Select([]string{"client_spm1_code", "client_type"}).
				Where("client_spm2_code = ?", p.Code2).
				First(&client).Error
			if err != nil {
				return err
			}
			// 校验 phone 是否被使用了, 被使用的话, 直接绑定到用户上
			err = db.Def().Select([]string{"id", "phone", "account", "open_id"}).Where("phone = ?", phone).First(&user).Error
			if err != nil {
				if !gorm.IsRecordNotFoundError(err) {
					return err
				}
				user.Phone = &phone
				user.OpenID = strings.ToUpper(randc.UUID())
				user.UserType = client.ClientType
				user.UserCode = uuid.GenerateForUIMS().String()
				user.NaCode = naCode
				err = db.Def().Save(&user).Error
				if err != nil {
					return err
				}
				var userInfo model.UserInfo
				userInfo.UserID = user.ID
				userInfo.NaCode = naCode
				userInfo.UserType = user.UserType
				userInfo.UserCode = user.UserCode
				userInfo.IsIdentify = "N"
				userInfo.Phone = phone
				err = db.Def().Save(&userInfo).Error
				if err != nil {
					return err
				}
			}
			// 查询到了 user 数据, 将 user_wechat 绑定到 user 上
			// 首先移除已绑定的 WeChatOpenId 数据
			db.Def().Where("wechat_open_id = ?", claims.OpenId).Delete(&model.UserWeChat{})
			// 然后绑定当前数据
			var userWeChat = model.UserWeChat{
				UserId:        uint(user.ID),
				WeChatId:      claims.WeChatId,
				Nickname:      req.Nickname,
				Sex:           req.Sex,
				Country:       req.Country,
				City:          req.City,
				Avatar:        req.Avatar,
				Province:      req.Province,
				WeChatOpenId:  claims.OpenId,
				WeChatUnionId: claims.UnionId,
			}
			err = db.Def().Save(&userWeChat).Error
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		// 进行登录
		userJwtService := jwtauth.UserJwtAuth{
			OpenId:   user.OpenID,
			ClientId: claims.ClientId,
			Account:  user.Account,
			State:    claims.State,
		}
		if userJwtService.IsFreeze() {
			return jwtauth.FreezeErr
		}
		authCode = userJwtService.GenerateCode()
		return nil
	}()
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Success(c, "", contexts.AppletInfoResp{
		Code:  authCode,
		State: state,
	})
	return
}
