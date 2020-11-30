package service

import (
	"math"
	"strconv"
	captcha_req "uims/internal/controllers/verifycode_controller/requests"
	captcha_resp "uims/internal/controllers/verifycode_controller/responses"
	"uims/pkg/tool"
)

// GenerateSlideCaptchaXYResponse 生成滑动式验证码定位图像所需要的坐标点以及本次数据的ID
func GenerateSlideCaptchaXYResponse(request *captcha_req.RequestGetSlideCaptchaCoordinate, resp *captcha_resp.SlideCapchaXYResponse) (err error) {
	defer func(e *error) {
		if er := recover(); er != nil {
			*e = er.(error)
		}
	}(&err)

	resp.SetX(tool.RandIntInRange(request.Xmin, request.Xmax))
	resp.SetY(tool.RandIntInRange(request.Ymin, request.Ymax))
	resp.ID = tool.GenerateRandStrWithCrypto(32)

	err = nil

	return
}

const ALPHABET = "abcdefghijklmnopqrstuvwxyz"

func MapSlideCaptchaData(value string) (result []map[string]int64) {
	valueBytes := []byte(value)
	result = make([]map[string]int64, len(valueBytes))
	for i, v := range valueBytes {
		result[i] = map[string]int64{
			string(ALPHABET[i]): int64(math.Pow(float64(v), 2)),
		}
	}
	return
}

func UnMapSlideCaptchaData(d []map[string]int) (x int, y int, err error) {
	defer func(e *error) {
		if er := recover(); er != nil {
			*e = er.(error)
		}
	}(&err)

	x, y = unMapSlideCaptchaData(d)

	return x, y, nil
}

// UnMapSlideCaptchaData 将前端传过来的坐标数据转换为x，y型的坐标
func unMapSlideCaptchaData(d []map[string]int) (int, int) {
	liststr := ""
	sepChar := ":"
	sepCharIndex := 0

	for i, v := range d {
		v1 := string(uint8(math.Sqrt(float64(v[string(ALPHABET[i])]))))
		if v1 == sepChar {
			sepCharIndex = i
		}
		liststr = liststr + v1
	}
	//fmt.Println(liststr)
	xint, err := strconv.Atoi(liststr[sepCharIndex+1:])
	if err != nil {
		panic(err)
	}
	xint = int(math.Sqrt(float64(xint)))

	yint1, err := strconv.Atoi(liststr[0:sepCharIndex])
	if err != nil {
		panic(err)
	}
	yint2, err := strconv.Atoi(liststr[sepCharIndex+1:])
	if err != nil {
		panic(err)
	}
	yint := int(math.Sqrt(float64(yint1 - yint2)))

	return xint, yint
}

func VerifyTheDisInAllowRange(x0, y0, x, y, allowDiff int) bool {
	if y0 == y && int(math.Abs(float64(x-x0))) <= allowDiff {
		return true
	}
	return false
}
