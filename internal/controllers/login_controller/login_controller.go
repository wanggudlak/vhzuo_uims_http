package login_controller

import (
	"errors"
	"fmt"
	"gitee.com/skysharing/vzhuo-go-captcha/captcha_api/math"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"time"
	"uims/conf"
	"uims/internal/controllers/auth_controller"
	"uims/internal/controllers/login_controller/contexts"
	resp "uims/internal/controllers/responses"
	"uims/internal/controllers/sms_controller"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/adapter"
	"uims/internal/service/jwtauth"
	"uims/pkg/db"
	"uims/pkg/encryption"
	"uims/pkg/gjwt"
	"uims/pkg/glog"
	"uims/pkg/jwt"
	"uims/pkg/tool"
)

// @Summary 用户登录鉴权
// @Produce  json
// @Param user body AuthenticateRequest
// @Success 200 {object} AuthenticateResp
// @Router /api/login/authenticate [post]
//
// 流程如下
// 1. 客户端web跳转到uims用户登录页
// 2. 用户提交登录表单
// 3. 登录凭据验证通过, 则 302 跳转到客户端web配置的回调地址中, 并带上 code + state
// 4. 客户端web获取到 code + state, 将请求客户端api, 将 code 换成 access_token
// 5. 客户端获取到 access_token 后,  通过 获取用户信息接口, 来查询用户基本信息
func Authenticate(c *gin.Context) {
	var req contexts.BaseAuthReq
	err := adapter.New().ReadRequestBody(c.Request, &req)
	if err != nil {
		fmt.Printf("adapter GetRequest err: %s\n", err.Error())
		resp.Error(c, err)
		return
	}

	switch req.AuthScene {
	case contexts.AccountPasswdVerifyCodeAuth: // 1用账号、密码、图片验证码请求登录
		AuthenticateByAccountPasswdVerifyCode(c)
		return
	case contexts.AccountPasswdVerifyCodeSMSCodeAuth: // 2用账号、密码、图片验证码、手机验证码请求登录
		AuthenticateByAccountPasswdVerifyCodeSMSCode(c)
		return
	case contexts.PhonePasswdAuth: // 3用手机号、密码请求登录
		AuthenticateByPhonePasswd(c)
		return
	case contexts.EmailPasswdAuth: // 4用邮箱、密码请求登录
		AuthenticateByEmailPasswd(c)
		return
	case contexts.PhoneVerifyCodeAndSlideCodeAuth: // 5,用手机号、滑动验证码+短信验证码请求登录
		AuthenticateByPhoneVerifyCodeAndSlideCode(c)
		return
	default:
		resp.Error(c, errors.New("未知的登录场景，请联系管理员进行配置"))
		return
	}
}

// AuthenticateByAccountPasswdVerifyCode 用账号、密码、图片验证码请求登录
func AuthenticateByAccountPasswdVerifyCode(c *gin.Context) {
	var err error
	req := contexts.AccountPasswdVerifyCodeAuthRequest{}
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}
	// 校验图片验证码
	if !math.Verify(req.VerificationKey, req.Code) && conf.Switch.ImgCaptcha {
		resp.Failed(c, "图片验证码输入错误", nil)
		return
	}

	user := model.User{}
	err = db.Def().
		Select([]string{"id", "open_id", "passwd", "passwd2", "encrypt_type", "salt", "status"}).
		Where(&model.User{Account: req.Account}).
		Where("isdel = ?", "N").
		First(&user).
		Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			resp.Failed(c, "用户不存在", nil)
		} else {
			resp.Error(c, err)
		}
		return
	}
	// 校验密码
	if err := checkLoginPasswd(req.AuthScene, req.Password, &user); err != nil {
		resp.Error(c, err)
		return
	}

	if err = checkUserStatus(user.Status); err != nil {
		resp.Error(c, err)
		return
	}

	// 生成授权码
	userJwtService := jwtauth.UserJwtAuth{
		OpenId:   user.OpenID,
		ClientId: uint(req.ClientId),
		Account:  req.Account,
		State:    req.State,
	}
	if userJwtService.IsFreeze() {
		resp.Error(c, jwtauth.FreezeErr)
		return
	}
	authCode := userJwtService.GenerateCode()

	// 构建鉴权响应
	authResp := &contexts.AuthenticateResp{}
	authResp, err = contexts.ParseAndPrepareAuthResponse(req.RedirectUrl, authCode, req.State)
	if err != nil {
		resp.Error(c, err)
		return
	}
	resp.Success(c, "success", authResp)
	return
}

// AuthenticateByAccountPasswdVerifyCodeSMSCode 用账号密码、图片验证码、手机短信验证码登录
// 主要分为两个流程
// 1. 用户直接提交账号密码, 这种情况下会校验账号密码并返回 password_token 作为下一步的凭据
// 2. 用户通过 /api/sms/verifycode/send 接口发送短信验证码
// 3. 再次调用登录接口, 带上短信验证码来进行校验, 校验通过, 即颁发 code
func AuthenticateByAccountPasswdVerifyCodeSMSCode(c *gin.Context) {
	var err error
	req := contexts.AccountPasswdVerifyCodeSMSCodeAuthRequest{}
	if err = c.ShouldBind(&req); err != nil {
		fmt.Printf("err2 : %+v \n ", err)
		resp.BadReq(c, err)
		return
	}

	type PasswordToken struct {
		gjwt.Jwt
		Account string
		Phone   string
	}

	// 第一步校验账号密码图片验证码成功后生成的passwd token
	// 校验 token 与短信验证码
	if len(req.PasswordToken) == 0 {
		err = func() error {
			if req.Account == "" {
				return fmt.Errorf("账号必须传入")
			}
			if req.Password == "" {
				return fmt.Errorf("密码必须传入")
			}
			// 校验图片验证码
			if !math.Verify(req.VerificationKey, req.Code) && conf.Switch.ImgCaptcha {
				return fmt.Errorf("图片验证码输入错误")
			}
			return nil
		}()
		if err != nil {
			resp.Error(c, err)
			return
		}

		user := model.User{}
		err = db.Def().
			Select([]string{"id", "open_id", "phone", "passwd", "passwd2", "encrypt_type", "salt", "status"}).
			Where(&model.User{Account: req.Account}).
			Where("isdel = ?", "N").
			First(&user).
			Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				resp.Failed(c, "用户不存在", nil)
			} else {
				resp.Error(c, err)
			}
			return
		}
		// 校验密码
		if err := checkLoginPasswd(req.AuthScene, req.Password, &user); err != nil {
			resp.Error(c, err)
			return
		}

		if err = checkUserStatus(user.Status); err != nil {
			resp.Error(c, err)
			return
		}

		claims := PasswordToken{
			Account: req.Account,
			Phone:   *user.Phone,
		}
		claims.SetAudience("password_token")
		claims.SetIssue()
		claims.SetTTL(5 * time.Minute)
		token, err := gjwt.CreateToken(&claims)
		if err != nil {
			resp.Error(c, err)
			return
		}
		resp.Success(c, "需要进行手机验证", gin.H{
			"password_token": token,
		})
		return
	} else {
		account, err := func() (string, error) {
			if req.SMSCode == "" {
				return "", fmt.Errorf("短信验证码必须传入")
			}
			claims := PasswordToken{}
			err = gjwt.Parse(req.PasswordToken, &claims)
			if err != nil {
				return "", fmt.Errorf("无效的 password_token")
			}
			scene := sms_controller.UseAccountFormat
			//fmt.Println("缓存key：", service.MakeSMSRedisCacheKey(claims.Phone, req.SPMFullCode+scene))
			//cacheCode := service.RedisDefaultCache.RedisGetStr(service.RedisKey("sms:" + req.SPMFullCode + ":" + claims.Phone))
			cacheCode := service.RedisDefaultCache.RedisGetStr(service.MakeSMSRedisCacheKey(claims.Phone, req.SPMFullCode+scene))
			//fmt.Println("输入验证码：", req.SMSCode)
			//fmt.Println("缓存验证码", cacheCode)
			if conf.Switch.SMSCaptcha && cacheCode != req.SMSCode {
				return "", fmt.Errorf("短信验证码错误")
			}
			return claims.Account, nil
		}()
		if err != nil {
			resp.Error(c, err)
			return
		}

		user := model.User{}
		err = db.Def().
			Select([]string{"id", "open_id", "phone", "passwd", "encrypt_type", "salt"}).
			Where(&model.User{Account: account}).
			First(&user).
			Error
		if err != nil {
			resp.Error(c, err)
			return
		}

		userJwtService := jwtauth.UserJwtAuth{
			OpenId:   user.OpenID,
			ClientId: uint(req.ClientId),
			Account:  user.Account,
			State:    req.State,
		}
		if userJwtService.IsFreeze() {
			resp.Error(c, jwtauth.FreezeErr)
			return
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
}

// AuthenticateByPhonePasswd
// 用手机号、密码请求登录
func AuthenticateByPhonePasswd(c *gin.Context) {
	var err error
	req := contexts.PhonePasswdAuthRequest{}
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}

	// 如果请求数据中的密码加密了，先解密
	err = (&req).DecryPasswdIfEncrypted()
	if err != nil {
		glog.Default().WithFields(glog.SetLogFiels(&map[string]interface{}{
			"err": err,
		})).Println("解密")
		resp.Error(c, err)
		return
	}

	// 查询用户信息
	user := model.User{}
	err = db.Def().
		Select([]string{"open_id", "status", "passwd", "passwd2", "salt", "encrypt_type"}).
		Where("phone = ?", req.Phone).
		Where("isdel = ?", "N").
		First(&user).
		Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			resp.Failed(c, "用户不存在", nil)
		} else {
			resp.Error(c, err)
		}
		return
	}

	// 校验密码
	if err := checkLoginPasswd(req.AuthScene, req.Password, &user); err != nil {
		resp.Failed(c, "密码错误", nil)
		return
	}

	if err = checkUserStatus(user.Status); err != nil {
		resp.Error(c, err)
		return
	}

	// 生成auth code
	userJwtService := jwtauth.UserJwtAuth{
		OpenId:   user.OpenID,
		ClientId: uint(req.ClientId),
		Account:  req.Phone,
		State:    req.State,
	}
	if userJwtService.IsFreeze() {
		resp.Error(c, jwtauth.FreezeErr)
		return
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

// AuthenticateByEmailPasswd
// 用邮箱、密码请求登录
func AuthenticateByEmailPasswd(c *gin.Context) {
	var err error
	req := contexts.EmailPasswdAuthRequest{}
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}

	// 查询用户信息
	user := model.User{}
	err = db.Def().
		Where("email = ?", req.Email).
		Where("isdel = ?", "N").
		First(&user).
		Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			resp.Failed(c, "用户不存在", nil)
		} else {
			resp.Error(c, err)
		}
		return
	}

	// 校验密码
	if err := checkLoginPasswd(req.AuthScene, req.Password, &user); err != nil {
		resp.Failed(c, "密码错误", nil)
		return
	}

	if err = checkUserStatus(user.Status); err != nil {
		resp.Error(c, err)
		return
	}

	// 生成auth code
	userJwtService := jwtauth.UserJwtAuth{
		OpenId:   user.OpenID,
		ClientId: uint(req.ClientId),
		Account:  req.Email,
		State:    req.State,
	}
	if userJwtService.IsFreeze() {
		resp.Error(c, jwtauth.FreezeErr)
		return
	}
	authCode := userJwtService.GenerateCode()
	authResp := &contexts.AuthenticateResp{}
	authResp, err = contexts.ParseAndPrepareAuthResponse(req.RedirectUrl, authCode, req.State)
	if err != nil {
		resp.Error(c, err)
		return
	}

	resp.Success(c, "success", authResp)
	//c.Redirect(http.StatusMovedPermanently, authResp.RedirectUrl)

	//c.Header("Cache-Control", "must-revalidate, no-store")
	//c.Header("Content-Type", " text/html;charset=UTF-8")
	//c.Redirect(http.StatusPermanentRedirect, authResp.RedirectUrl)

	//c.Abort()
	return
}

// AuthenticateByPhoneVerifyCode
// 用手机号、手机短信验证码请求登录，获取手机短信验证码之前需要校验滑动验证码
// 获取短信验证码之前通过滑动验证码校验通过之后才发送短信
// 如果发现未注册，先走注册的逻辑然后继续走登录
func AuthenticateByPhoneVerifyCodeAndSlideCode(c *gin.Context) {
	var err error
	req := contexts.PhoneVerifyCodeAndSlideCodeAuthRequest{}
	if err = c.ShouldBind(&req); err != nil {
		resp.BadReq(c, err)
		return
	}

	// 校验短信验证码
	scene := sms_controller.UsePhoneLoginRegister
	//cacheCode := service.RedisDefaultCache.RedisGetStr(service.RedisKey("sms:" + req.SPMFullCode + ":" + req.Phone))
	cacheCode := service.RedisDefaultCache.RedisGetStr(service.MakeSMSRedisCacheKey(req.Phone, req.VerificationKey+scene))
	if conf.Switch.SMSCaptcha && cacheCode != req.SMSCode {
		resp.Failed(c, "短信验证码错误", nil)
		return
	}

	// 查询用户信息
	pUser := &model.User{}
	err = db.Def().
		Where("phone = ?", req.Phone).
		Where("isdel = ?", "N").
		First(pUser).
		Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			//resp.Failed(c, "用户不存在", nil)
			// 走注册，然后继续
			onlyRegisterReq := auth_controller.RegisterRequest{
				Phone:       req.Phone,
				SMSCode:     req.SMSCode,
				Password:    string(tool.GenerateRandBytesWithCrypto(6)),
				SPMFullCode: req.SPMFullCode,
			}
			pUser, err = auth_controller.DoRegister(&onlyRegisterReq)
			if err != nil {
				resp.Error(c, err)
				return
			}
		} else {
			resp.Error(c, err)
		}
	}

	if err = checkUserStatus(pUser.Status); err != nil {
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
	if userJwtService.IsFreeze() {
		resp.Error(c, jwtauth.FreezeErr)
		return
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

// checkLoginPasswd
// 说明：由于结算系统和任务系统之前采用的登录密码生成算法不同，本次开发uims系统，暂需要设置两个密码保存字段
// 任务系统用passwd2校验，其它用passwd校验；如果用户申请更新密码，我们会统一密码生成为一致的，即统一采用目前结算系统所用之算法
func checkLoginPasswd(authScene int, reqPasswd string, user *model.User) error {
	switch user.EncryptType {
	default:
		fallthrough
	case 0:
		// 默认鉴权 / 结算系统鉴权
		if !encryption.BcryptCheck(reqPasswd, user.Passwd) {
			return errors.New("密码错误")
		}
	case 1:
		if authScene == contexts.AccountPasswdVerifyCodeAuth || authScene == contexts.AccountPasswdVerifyCodeSMSCodeAuth {
			if !encryption.BcryptCheck(reqPasswd, user.Passwd) {
				return errors.New("密码错误")
			} else {
				return nil
			}
		} else {
			// 任务系统鉴权
			//fmt.Printf("数据库密码是: %s", user.Passwd2)
			if !encryption.DefaultPBKDF2Options.CheckPBKDF2PasswdForVzhuoTaskSYS(reqPasswd, user.Passwd2, user.Salt) {
				return errors.New("密码错误")
			}
		}
	}
	return nil
}

type LoginResult struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

//后台管理员登陆
func BackgroundLogin(c *gin.Context) {
	var backgroundLoginRequest contexts.BackgroundLoginRequest
	var err error
	if err = c.ShouldBindJSON(&backgroundLoginRequest); err != nil {
		resp.BadReq(c, err)
		return
	}
	user := model.User{}
	err = service.GetUserService().GetUserInfoByAccount(&user, backgroundLoginRequest.Account)
	if err != nil {
		resp.Failed(c, "服务器错误", nil)
		return
	}
	if user.ID == 0 {
		resp.Failed(c, "账号错误", nil)
		return
	}

	if user.Status == "N" || user.Isdel == "Y" {
		resp.Failed(c, "该账号暂不可用", nil)
		return
	}

	if user.UserType != "UIMS" {
		resp.Failed(c, "无权限登陆系统", nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(backgroundLoginRequest.Passwd)); err != nil {
		resp.Failed(c, "密码错误", nil)
		return
	}

	//缓存用户信息
	err = service.CacheBackgroundUserInfo(&user, user.Account)
	if err != nil {
		resp.Error(c, err)
		return
	}
	generateToken := GenerateToken(user)
	//获取token解析后的数据
	//claims := c.MustGet("claims").(*jwt.CustomClaims)
	resp.Success(c, "success", generateToken)
}

// 生成令牌  创建jwt风格的token
func GenerateToken(user model.User) LoginResult {
	j := &jwt.JWT{
		[]byte("UIMS-BACK-JWT"),
	}
	claims := jwt.CustomClaims{
		user.Account,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),   // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600*6), // 过期时间 六小时
			Issuer:    "UIMS-BACK",                       //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		return LoginResult{
			User:  user,
			Token: token,
		}
	}

	data := LoginResult{
		User:  user,
		Token: token,
	}

	return data
}

func GetRSAPubKey(c *gin.Context) {
	resp.Success(c, "success", gin.H{"k": string(encryption.GetAPPPbulicKeyContent())})
}

// 检查用户冻结状态
func checkUserStatus(status string) error {
	if status == "N" {
		return fmt.Errorf("用户已冻结, 无法登录")
	}
	return nil
}
