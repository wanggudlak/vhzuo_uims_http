package sms_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
	"uims/conf"
	responses2 "uims/internal/controllers/responses"
	requests2 "uims/internal/controllers/sms_controller/requests"
	"uims/internal/controllers/sms_controller/responses"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/email"
	"uims/pkg/db"
	"uims/pkg/glog"
	"uims/pkg/tool"

	"gitee.com/skysharing/vzhuo-go-sms/smsapi"
	"github.com/gin-gonic/gin"
)

const (
	UsePhoneLogin         = "phone:login"
	UsePhoneRegister      = "phone:register"
	UsePhoneRegisterLogin = "phone:register:login" // 用手机号注册，注册后直接登录
	UsePhoneLoginRegister = "phone:login:register" // 用手机号登录，如果用户未注册，直接注册后登录
	UsePhoneFindPasswd    = "phone:findpasswd"
	UseEmailLogin         = "email:login"
	UseEmailRegister      = "email:register"
	UseEmailFindPasswd    = "email:findpasswd"
	UsePhoneBindWeChat    = "phone:bind:wechat" // 任务系统绑定微信
	UseAccountFormat      = "account:format"
)

var (
	smsConf       = conf.SMS
	cacheDuration = 5 * time.Minute
)

// @Summary 发送短信验证码
// @Produce  json
// @Param user body requests.SMSgetVerifyCodeRequest
// @Success 200 {object} smsresp.SMSaliResponse
// @Router /api/sms/verifycode/send [get]
// Send 发送短信验证码
func Send(c *gin.Context) {
	scene := c.DefaultQuery("scene", "")
	switch scene {
	case UsePhoneLogin, UsePhoneRegister, UsePhoneLoginRegister, UsePhoneRegisterLogin, UsePhoneFindPasswd, UsePhoneBindWeChat:
		req := requests2.SMSCodeByPhoneRequest{}
		if err := c.ShouldBind(&req); err != nil {
			responses2.BadReq(c, err)
			return
		}
		can, err := CanSendSMS(scene, req.Phone)
		if err != nil {
			responses2.Error(c, err)
			return
		}
		if !can {
			responses2.Failed(c, "暂不能发送验证码，请联系管理员", nil)
			return
		}
		//if UsePhoneLogin == scene {
		//	// 判断用户是否已经注册
		//	user := model.User{}
		//	err := db.Def().
		//		Select([]string{"id"}).
		//		Where(&model.User{Phone: req.Phone}).
		//		Where("isdel = ?", "N").
		//		Where("status = ?", "Y").
		//		First(&user).
		//		Error
		//	if err != nil {
		//		if gorm.IsRecordNotFoundError(err) {
		//			responses2.Failed(c, "用户不存在", nil)
		//		} else {
		//			responses2.Error(c, err)
		//		}
		//		return
		//	}
		//}

		err = SendPhoneSMSCodeAPI(req.Phone, req.SPM, scene)
		if err != nil {
			glog.Default().Printf("通过手机发送短信验证码后失败: %s\n", err.Error())
			responses2.Failed(c, fmt.Sprintf("发送短信验证码失败：%s", err.Error()), nil)
			return
		}
		responses2.Success(c, "success", nil)
		return
	case UseEmailLogin, UseEmailRegister, UseEmailFindPasswd:
		req := requests2.SMSCodeByEmailRequest{}
		if err := c.ShouldBind(&req); err != nil {
			responses2.BadReq(c, err)
			return
		}

		can, err := CanSendSMS(scene, req.Email)
		if err != nil {
			responses2.Error(c, err)
			return
		}
		if !can {
			responses2.Failed(c, "暂不能发送验证码，请联系管理员", nil)
			return
		}

		//if UseEmailLogin == scene {
		//	// 判断用户是否已经注册
		//	user := model.User{}
		//	err := db.Def().
		//		Select([]string{"id"}).
		//		Where(&model.User{Email: req.Email}).
		//		Where("isdel = ?", "N").
		//		Where("status = ?", "Y").
		//		First(&user).
		//		Error
		//	if err != nil {
		//		if gorm.IsRecordNotFoundError(err) {
		//			responses2.Failed(c, "用户不存在", nil)
		//		} else {
		//			responses2.Error(c, err)
		//		}
		//		return
		//	}
		//}

		randIntStr := tool.GenerateRandStrWithMath(6, []byte("123456789"))
		bodyFormat := "【微桌】验证码%s，用于微桌系统登录，该验证码%d分钟内有效。工作人员不会向您索要验证码，请勿泄漏于他人，以免造成损失！"
		body := fmt.Sprintf(bodyFormat, randIntStr, int(cacheDuration/time.Minute))
		context := &email.Context{
			To:       []string{req.Email},
			Subject:  "微桌更新密码",
			BodyType: email.HTMlContentType,
			Body:     body,
		}
		err = context.Send()
		if err != nil {
			glog.Default().Printf("通过邮箱发送短信验证码失败: %s\n", err.Error())
			responses2.Failed(c, fmt.Sprintf("发送短信验证码失败：%s", err.Error()), nil)
			return
		}
		err = service.RedisDefaultCache.RedisCacheString(service.MakeSMSRedisCacheKey(req.Email, req.SPM+scene), randIntStr, cacheDuration)
		if err != nil {
			glog.Default().Printf("通过邮箱发送短信验证码后缓存失败: %s\n", err.Error())
			responses2.Failed(c, "发送短信验证码失败", nil)
			return
		}
		responses2.Success(c, "success", nil)
		return
	case UseAccountFormat:
		req := requests2.SMSCodeByAccountFormatRequest{}
		if err := c.ShouldBind(&req); err != nil {
			responses2.BadReq(c, err)
			return
		}
		//如果是邮箱
		if tool.VerifyEmailFormat(req.Account) {
			err := SendEmailSMSCode(req.Scene, req.Account, req.SPM, req.LoginAccount)
			if err != nil {
				responses2.Failed(c, err.Error(), nil)
				return
			}
			responses2.Success(c, "success", nil)
			return
		} else {
			//手机号
			err := SendPhoneSMSCode(req.Scene, req.Account, req.SPM)
			if err != nil {
				glog.Default().Printf("通过手机发送短信验证码后失败: %s\n", err.Error())
				responses2.Failed(c, fmt.Sprintf("发送短信验证码失败：%s", err.Error()), nil)
				return
			}

			responses2.Success(c, "success", nil)
		}
		return
	default:
		responses2.Failed(c, "需要指明使用什么媒介发送信息", nil)
		return
	}
}

// phone
// k 关键字
// scene 场景
func SendPhoneSMSCodeAPI(phone string, k string, scene string) error {
	randIntStr := tool.GenerateRandStrWithMath(6, []byte("123456789"))
	glog.Channel("sms").WithFields(log.Fields{
		"phone": phone,
		"code":  randIntStr,
		"k":     k,
		"scene": scene,
	}).Info("记录短信发送")
	smsDriverType := smsConf.GetDriverType()
	if len(smsDriverType) == 0 {
		return errors.New("SMS driver is invalid.")
	}
	switch smsDriverType {
	case conf.SMS_ALI_DRIVER:
		smsAliConf := smsConf.GetAliDriverParam()
		//smsAliRegionID := smsAliConf.GetSMSaliRegionID()
		client, err := smsapi.NewAliSMSclient(smsAliConf.GetSMSaliRegionID(), smsAliConf.GetSMSaliAccessKey(),
			smsAliConf.GetSMSaliAccessSecret())
		if err != nil {
			glog.Channel("sms").Printf("通过手机发送短信验证码失败: %s\n", err.Error())
			return err
		}
		smsRequest := smsapi.NewAliSMSrequest(smsAliConf.GetSMSaliRegionID(), phone, conf.SMS_ALI_SIGN,
			conf.SMS_ALI_VERIFY_CODE_TEMPLATE_ID, "{\"code\":"+randIntStr+"}")
		smsResponse, err := client.SendSMS(smsRequest)
		if err != nil {
			glog.Channel("sms").Printf("通过手机发送短信验证码失败: %s\n", err.Error())
			return err
		}
		content := &responses.SMSaliResponse{}
		if err := json.Unmarshal([]byte(smsResponse.GetHttpContentString()), content); err != nil {
			return err
		}
		if smsResponse.IsSuccess() {
			isSendSuccess, message := content.IsSendSuccess()
			if isSendSuccess {
				err = service.RedisDefaultCache.RedisCacheString(service.MakeSMSRedisCacheKey(phone, k+scene), randIntStr, cacheDuration)
				if err != nil {
					glog.Channel("sms").Printf("通过手机发送短信验证码后缓存失败: %s\n", err.Error())
					return err
				}

				return nil
			} else {
				return errors.New(message)
			}
		} else {
			return errors.New("Send SMS failed.")
		}
	default:
		return errors.New("SMS driver is invalid")
	}
}

// @Summary 验证验证码
// @Produce  json
// @Param user body requests.VerifyRequest
// @Router /api/sms/verifycode/verify [post]
// Verify 验证验证码
func Verify(c *gin.Context) {
	req := requests2.VerifyRequest{}
	if err := c.ShouldBind(&req); err != nil {
		responses2.BadReq(c, err)
		return
	}

	//req.By = strings.TrimSpace(req.By)
	//req.Code = strings.TrimSpace(req.Code)

	cacheK := service.MakeSMSRedisCacheKey(req.By, req.Key)
	codeInCache := service.RedisDefaultCache.RedisGetStr(cacheK)
	if req.Code != codeInCache {
		responses2.Failed(c, "验证码错误", nil)
		return
	}
	service.RedisDefaultCache.RedisDel(cacheK)
	responses2.Success(c, "success", nil)
	return
}

// CanSendSMS 判断各种业务场景下能否发送验证码
// 目前存在的业务场景：
//UsePhoneLogin      = "phone:login"       用手机号登录
//UsePhoneRegister   = "phone:register"    用手机号注册
//UsePhoneFindPasswd = "phone:findpasswd"  用手机号找回密码
//UseEmailLogin      = "email:login"       用邮箱登录
//UseEmailRegister   = "email:register"    用邮箱注册
//UseEmailFindPasswd = "email:findpasswd"  用邮箱找回密码
func CanSendSMS(scene string, fieldV interface{}) (bool, error) {
	switch scene {
	case UsePhoneRegister, UsePhoneRegisterLogin:
		exist, err := service.IsExistByPhone(fieldV.(string))
		if err != nil {
			return false, err
		} else {
			if exist {
				return false, errors.New("用户已经注册")
			} else {
				return true, nil
			}
		}
	case UsePhoneLogin, UsePhoneFindPasswd:
		exist, err := service.IsExistByPhone(fieldV.(string))
		if err != nil {
			return false, err
		} else {
			if !exist {
				return false, errors.New("用户未注册")
			} else {
				return true, nil
			}
		}
	case UsePhoneLoginRegister:
		return true, nil
	case UseEmailRegister:
		exist, err := service.IsExistByEmail(fieldV.(string))
		if err != nil {
			return false, err
		} else {
			if exist {
				return false, errors.New("用户已经注册")
			} else {
				return true, nil
			}
		}
	case UseEmailLogin, UseEmailFindPasswd:
		exist, err := service.IsExistByEmail(fieldV.(string))
		if err != nil {
			return false, err
		} else {
			if !exist {
				return false, errors.New("用户未注册")
			} else {
				return true, nil
			}
		}
	case UsePhoneBindWeChat:
		break
	default:
		return false, errors.New("未知的业务场景，不能发送验证码")
	}
	return true, nil
}

func IsExistUserByScene(scene string, fieldV interface{}) (bool, error) {
	switch scene {
	case UsePhoneLogin, UsePhoneFindPasswd, UsePhoneRegister, UsePhoneRegisterLogin:
		return service.IsExistByPhone(fieldV.(string))
	case UseEmailLogin, UseEmailRegister, UseEmailFindPasswd:
		return service.IsExistByEmail(fieldV.(string))
	default:
		return false, errors.New("未知的查询条件，无法查询用户")
	}
}

func SendPhoneSMSCode(scene string, loginPhone string, spm string) error {
	can, err := CanSendSMS(UsePhoneLogin, loginPhone)
	if err != nil {
		return err
	}
	if !can {
		return errors.New("暂不能发送验证码，请联系管理员")
	}
	err = SendPhoneSMSCodeAPI(loginPhone, spm, scene)
	if err != nil {
		glog.Default().Printf("通过手机发送短信验证码后失败: %s\n", err.Error())
		return errors.New(fmt.Sprintf("发送短信验证码失败：%s", err.Error()))
	}
	return nil
}

func SendEmailSMSCode(scene string, loginEmail string, spm string, loginAccount string) error {
	can, err := CanSendSMS(UseEmailLogin, loginEmail)
	if err != nil {
		return err
	}
	if !can {
		return errors.New("暂不能发送验证码，请联系管理员")
	}

	var userAuth model.User
	err = db.Def().Where("account = ?", loginAccount).First(&userAuth).Error
	if err != nil || userAuth.ID == 0 {
		return errors.New("账号信息错误")
	}
	randIntStr := tool.GenerateRandStrWithMath(6, []byte("123456789"))
	bodyFormat := "【微桌】验证码%s，用于微桌系统登录，该验证码%d分钟内有效。工作人员不会向您索要验证码，请勿泄漏于他人，以免造成损失！"
	body := fmt.Sprintf(bodyFormat, randIntStr, int(cacheDuration/time.Minute))
	context := &email.Context{
		To:       []string{loginEmail},
		Subject:  "微桌登陆密码",
		BodyType: email.HTMlContentType,
		Body:     body,
	}
	err = context.Send()
	if err != nil {
		glog.Default().Printf("通过邮箱发送短信验证码失败: %s\n", err.Error())
		return errors.New(fmt.Sprintf("发送短信验证码失败：%s", err.Error()))
	}
	fmt.Println("发送验证码缓存key：", service.MakeSMSRedisCacheKey(loginEmail, spm+scene))
	fmt.Println("发送验证码缓存value：", randIntStr)
	err = service.RedisDefaultCache.RedisCacheString(service.MakeSMSRedisCacheKey(*userAuth.Phone, spm+scene), randIntStr, cacheDuration)
	if err != nil {
		glog.Default().Printf("通过邮箱发送短信验证码后缓存失败: %s\n", err.Error())
		return errors.New("发送短信验证码失败")
	}
	return nil
}
