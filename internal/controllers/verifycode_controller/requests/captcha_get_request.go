package requests

const (
	DefaultHeight = 60
	DefaultWidth  = 100
)

type RequestOfGenerateMathCaptcha struct {
	Height int `json:"height" form:"height" binding:"" faker:"height" example:"60" comment:"验证码图片高度"`
	Width  int `json:"width" form:"width" binding:"" faker:"width" example:"100" comment:"验证码图片宽度"`
}

func (req *RequestOfGenerateMathCaptcha) SetDefaultSize() *RequestOfGenerateMathCaptcha {
	if req.Height == 0 {
		req.Height = DefaultHeight
	}
	if req.Width == 0 {
		req.Width = DefaultWidth
	}
	return req
}

type RequestGetSlideCaptchaCoordinate struct {
	Xmin int64 `json:"x_min" form:"x_min" binding:"" faker:"x_min" example:"80" comment:"X坐标最小值"`
	Xmax int64 `json:"x_max" form:"x_max" binding:"" faker:"x_max" example:"230" comment:"X坐标最大值"`
	Ymin int64 `json:"y_min" form:"y_min" binding:"" faker:"y_min" example:"40" comment:"Y坐标最小值"`
	Ymax int64 `json:"y_max" form:"y_max" binding:"" faker:"y_max" example:"130" comment:"Y坐标最大值"`
}

func (req *RequestGetSlideCaptchaCoordinate) SetDefaultXY() *RequestGetSlideCaptchaCoordinate {
	req.Xmin = 80
	req.Xmax = 230
	req.Ymin = 40
	req.Ymax = 130
	return req
}

type RequestVerifySlideCaptchaLocation struct {
	Phone      string           `json:"phone" form:"phone" binding:"required,len=11,mobile" example:"13517210606" comment:"手机号"`
	Scene      string           `json:"scene" form:"scene" binding:"required" faker:"scene" example:"phone" comment:"场景，用什么媒介发送短信"`
	CaptchaID  string           `json:"id" form:"id" binding:"required" faker:"id" example:"fjaksdjfkafads" comment:"验证码ID"`
	Coordinate []map[string]int `json:"xy" form:"xy" binding:"required" faker:"" example:"" comment:"验证码坐标位置数据"`
}
