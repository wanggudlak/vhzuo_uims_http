package auth_controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
	"uims/conf"
	"uims/internal/controllers/login_controller/contexts"
	resp "uims/internal/controllers/responses"
	"uims/internal/controllers/sms_controller"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/jwtauth"
	"uims/internal/service/uuid"
	"uims/pkg/db"
	"uims/pkg/encryption"
	"uims/pkg/gjwt"
	"uims/pkg/glog"
	"uims/pkg/randc"
)

type FindPasswordClaims struct {
	gjwt.Jwt
	UserId int
}

// 获取修改密码的 find_password_token
// 需要校验手机号/邮箱与验证码
func FindPasswordToken(c *gin.Context) {
	var err error
	var req FindPasswordTokenRequest
	if err = c.ShouldBind(&req); err != nil {
		glog.Default().Println("-1:find_password_token_err=", err.Error())
		resp.BadReq(c, err)
		return
	}
	account := ""
	where := ""
	scene := ""
	if len(req.Phone) != 0 {
		account = req.Phone
		where = "phone = ?"
		scene = sms_controller.UsePhoneFindPasswd
	} else if len(req.Email) != 0 {
		account = req.Email
		where = "email = ?"
		scene = sms_controller.UseEmailFindPasswd
	} else {
		resp.Failed(c, "手机号或邮箱是必需的", nil)
		return
	}

	var token string
	token, err = func() (string, error) {
		var user model.User
		cacheCode := service.RedisDefaultCache.RedisGetStr(service.MakeSMSRedisCacheKey(account, req.SPMFullCode+scene))
		if conf.Switch.SMSCaptcha && cacheCode != req.SMSCode {
			return "", errors.New("短信验证码错误")
		}
		err = db.Def().Where(where, account).First(&user).Error
		if err != nil {
			glog.Default().Println("0:find_password_token_err=", err.Error())
			return "", err
		}
		claims := &FindPasswordClaims{
			UserId: user.ID,
		}
		claims.SetIssue()
		claims.SetAudience("find_password")
		claims.SetTTL(5 * time.Minute) // 5分钟有效期
		token, err := gjwt.CreateToken(claims)
		if err != nil {
			glog.Default().Println("1:find_password_token_err=", err.Error())
			return "", err
		}
		return token, nil
	}()
	if err != nil {
		glog.Default().Println("2:find_password_token_err=", err.Error())
		resp.Error(c, err)
		return
	}

	resp.Success(c, "", gin.H{
		"find_password_token": token,
	})
	return
}

//// MakeTokenForFindPasswdForm 忘记密码表单获取标识表单唯一性的token
//func MakeTokenForFindPasswdForm(c *gin.Context) {
//	req := GetTokenForFindPasswdFormRequest{}
//	if err := c.ShouldBind(&req); err != nil {
//		resp.BadReq(c, err)
//		return
//	}
//	csrfToken, err := service.GenerateCSRFToken(req.ClientId)
//	if err != nil {
//		resp.Error(c, err)
//		return
//	}
//	resp.Success(c, "success", gin.H{
//		"_forget_passwd_token": csrfToken,
//	})
//}

// 传入 find_password_token 修改为新密码
func FindPassword(c *gin.Context) {
	var err error
	var req FindPasswordRequest
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	err = func() error {
		claims := FindPasswordClaims{}
		err = gjwt.Parse(req.FindPasswordToken, &claims)
		if err != nil {
			return err
		}
		if !claims.VerifyAudience("find_password", true) {
			return fmt.Errorf("token 场景不正确")
		}
		// 修改用户密码
		err = service.UserService{}.SetPassword(claims.UserId, req.NewPassword)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Success(c, "", nil)
	return
}

// 通过手机号验证码快速注册用户
// 手机号
// 验证码
// 密码
func Register(c *gin.Context) {
	var err error
	var req RegisterRequest
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	_, err = DoRegister(&req)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Success(c, "", nil)
	return
}

// RegisterAndLogin 实现注册后即登录
// 手机号
// 验证码
// 密码
func RegisterAndLogin(c *gin.Context) {
	var req RegisterAndLoginRequest
	if err := c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	onlyRegisterReq := RegisterRequest{
		Phone:       req.Phone,
		SMSCode:     req.SMSCode,
		Password:    req.Password,
		SPMFullCode: req.SPMFullCode,
	}
	pUser, err := DoRegister(&onlyRegisterReq)
	if err != nil {
		resp.Error(c, err)
		return
	}
	// 生成auth code
	userJwtService := jwtauth.UserJwtAuth{
		OpenId:   pUser.OpenID,
		ClientId: uint(req.ClientId),
		Account:  req.Phone,
		State:    req.State,
	}
	authCode := userJwtService.GenerateCode()
	authResp := &contexts.AuthenticateResp{}
	authResp, err = contexts.ParseAndPrepareAuthResponse(req.RedirectUrl, authCode, req.State)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Success(c, "success", authResp)
	return
}

func DoRegister(req *RegisterRequest) (pUser *model.User, err error) {
	pUser = &model.User{}
	err = db.Def().Transaction(func(tx *gorm.DB) error {
		err = db.Def().Where("phone = ?", req.Phone).First(pUser).Error
		if err == nil {
			return fmt.Errorf("手机号已被注册")
		} else if !gorm.IsRecordNotFoundError(err) {
			return err
		}
		p := service.ParseSPMstring(req.SPMFullCode)
		var client model.Client
		err = db.Def().Select([]string{"client_spm1_code", "client_type"}).
			Where("client_spm2_code = ?", p.Code2).
			First(&client).Error
		if err != nil {
			return err
		}
		pUser.Phone = &req.Phone
		pUser.OpenID = strings.ToUpper(randc.UUID())
		pUser.UserType = client.ClientType
		pUser.UserCode = uuid.GenerateForUIMS().String()
		pUser.NaCode = "+86"
		pUser.Passwd, err = encryption.BcryptHash(req.Password)
		pUser.EncryptType = 0
		if err != nil {
			return err
		}
		err = db.Def().Save(pUser).Error
		if err != nil {
			return err
		}
		var userInfo model.UserInfo
		userInfo.UserID = pUser.ID
		userInfo.NaCode = pUser.NaCode
		userInfo.Phone = *pUser.Phone
		userInfo.UserType = pUser.UserType
		userInfo.UserCode = pUser.UserCode
		userInfo.IsIdentify = "N"
		err = db.Def().Save(&userInfo).Error
		if err != nil {
			return err
		}
		return nil
	})

	return
}
