package verifycode_controller

import (
	"fmt"
	"gitee.com/skysharing/vzhuo-go-captcha/captcha_api/math"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
	responses2 "uims/internal/controllers/responses"
	"uims/internal/controllers/sms_controller"
	requests2 "uims/internal/controllers/verifycode_controller/requests"
	"uims/internal/controllers/verifycode_controller/responses"
	"uims/internal/service"
	"uims/pkg/glog"
)

// @Summary 生成数学验证码
// @Produce  json
// @Param user body requests.RequestOfGenerateMathCaptcha
// @Success 200 {object} verifyresp.CaptchaGetRequest
// @Router /api/captcha/math [get]
// GenerateMathCaptchaBase64
func GenerateMathCaptchaBase64(c *gin.Context) {
	var request requests2.RequestOfGenerateMathCaptcha
	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	request.SetDefaultSize()

	resp := responses.MathCaptchaResponse{}
	err := APIGenerateMathCaptchaBase64(&request, &resp)
	if err != nil {
		responses2.Error(c, err)
		return
	}

	responses2.Success(c, "生成数学图片验证码成功", resp)
	return
}

// APIGenerateMathCaptchaBase64 根据生成数学图片验证码的请求生成验证码并以响应形式返回
func APIGenerateMathCaptchaBase64(request *requests2.RequestOfGenerateMathCaptcha, resp *responses.MathCaptchaResponse) error {
	id, b64s, err := math.NewMathCaptchaDriver().
		SetRequiredParam(request.Height, request.Width).
		GenerateMathCaptcha()
	if err != nil {
		return nil
	}

	resp.ID = id
	resp.ImgB64 = b64s

	return nil
}

// @Summary 生成滑动式验证码所需要的坐标点及ID
// @Produce  json
// @Param user body requests.RequestGetSlideCaptchaCoordinate
// @Success 200 {object} verifyresp.SlideCapchaXYResponse
// @Router /api/captcha/slide [get]
// GenerateSlideRangeCoordinatePoints
func GenerateSlideRangeCoordinatePoints(c *gin.Context) {
	var request requests2.RequestGetSlideCaptchaCoordinate
	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	request.SetDefaultXY()

	resp := responses.SlideCapchaXYResponse{}
	err := service.GenerateSlideCaptchaXYResponse(&request, &resp)
	if err != nil {
		responses2.Error(c, err)
		return
	}

	// 本次生成的数据以resp.ID为key缓存至缓存系统
	xy := fmt.Sprintf("%d:%d", resp.GetX(), resp.GetY())
	err = service.RedisDefaultCache.RedisCacheString(resp.ID, xy, time.Minute*5)
	if err != nil {
		responses2.Error(c, err)
		return
	}

	// 将结果转换为指定数据结构后返回给前端
	resp.MapXY = service.MapSlideCaptchaData(xy)

	responses2.Success(c, "success", resp)
	return
}

// @Summary 验证滑动验证码位置是否正确，如果正确就发送短信验证码
// @Produce  json
// @Param user body requests.RequestVerifySlideCaptchaLocation
// @Success 200 {object}
// @Router /api/captcha/verifyslide [post]
// VerifySlideCaptchaLocationAndSendSMS
func VerifySlideCaptchaLocationAndSendSMS(c *gin.Context) {
	var request requests2.RequestVerifySlideCaptchaLocation
	if err := c.ShouldBindJSON(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	cacheData := service.RedisDefaultCache.RedisGetStr(request.CaptchaID)
	if len(cacheData) == 0 {
		responses2.Failed(c, "图片已失效,请刷新重试", nil)
		return
	}

	xyByte := strings.Split(cacheData, ":")
	if len(xyByte) != 2 {
		e := service.RedisDefaultCache.RedisDel(request.CaptchaID)
		if e != nil {
			glog.Default().Errorf("清除滑动式验证码缓存时失败：[id: %s]<%s>\n", request.CaptchaID, e.Error())
		}
		responses2.Failed(c, "图片已失效,请刷新重试", nil)
		return
	}
	x0, err := strconv.Atoi(xyByte[0])
	if err != nil {
		responses2.Failed(c, "验证失败,请刷新重试", nil)
		return
	}
	y0, err := strconv.Atoi(xyByte[1])
	if err != nil {
		responses2.Failed(c, "验证失败,请刷新重试", nil)
		return
	}

	x, y, err := service.UnMapSlideCaptchaData(request.Coordinate)
	if err != nil {
		responses2.Failed(c, "验证失败,请刷新重试", nil)
		return
	}

	if !service.VerifyTheDisInAllowRange(x0, y0, x, y, 5) {
		responses2.Failed(c, "验证失败,请刷新重试", nil)
		return
	}

	e := service.RedisDefaultCache.RedisDel(request.CaptchaID)
	if e != nil {
		glog.Default().Errorf("清除滑动式验证码缓存时失败：[id: %s]<%s>\n", request.CaptchaID, e.Error())
	}

	can, err := sms_controller.CanSendSMS(request.Scene, request.Phone)
	if err != nil {
		responses2.Error(c, err)
		return
	}
	if !can {
		responses2.Failed(c, "暂不能发送验证码，请联系管理员", nil)
		return
	}
	// 判断当前手机号是否注册
	// 查询用户信息
	//user := model.User{}
	//err = db.Def().
	//	Where("phone = ?", request.Phone).
	//	Where("isdel = ?", "N").
	//	Where("status = ?", "Y").
	//	First(&user).
	//	Error
	//if err != nil {
	//	if gorm.IsRecordNotFoundError(err) {
	//		responses2.Failed(c, "用户不存在", nil)
	//	} else {
	//		responses2.Error(c, err)
	//	}
	//	return
	//}
	// 给手机号发送短信验证码
	go func() {
		//err = sms_controller.SendPhoneSMSCodeAPI(request.Phone, "sms:"+request.Phone+request.CaptchaID)
		err = sms_controller.SendPhoneSMSCodeAPI(request.Phone, request.CaptchaID, request.Scene)
		if err != nil {
			glog.Default().WithField("phone", request).Errorf("发送手机验证码失败：%s", err.Error())
			//responses2.Failed(c, fmt.Sprint("发送手机验证码失败：%s", err.Error()), nil)
			//return
		} else {
			glog.Default().WithField("phone", request).Info("发送手机验证码成功")
		}
	}()

	responses2.Success(c, "验证成功！", nil)
}
