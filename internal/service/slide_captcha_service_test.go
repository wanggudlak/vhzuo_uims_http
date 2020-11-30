package service_test

import (
	"fmt"
	"testing"
	captcha_req "uims/internal/controllers/verifycode_controller/requests"
	captcha_resp "uims/internal/controllers/verifycode_controller/responses"
	"uims/internal/service"
)

func TestGenerateSlideCaptchaXYResponse(t *testing.T) {
	req := captcha_req.RequestGetSlideCaptchaCoordinate{}
	resp := captcha_resp.SlideCapchaXYResponse{}

	req.Xmin = 40
	req.Xmax = 230
	req.Ymin = 20
	req.Ymax = 250

	err := service.GenerateSlideCaptchaXYResponse(&req, &resp)

	if err != nil {
		t.Errorf("GenerateSlideCaptchaXYResponse error: <%s>\n", err.Error())
	}

	fmt.Println(resp)
}

func TestMapSlideCaptchaData(t *testing.T) {
	testV := "100:50"
	result := service.MapSlideCaptchaData(testV)
	fmt.Println(result)
}

func TestUnMapSlideCaptchaData(t *testing.T) {
	testV := []map[string]int{
		{"a": 2916},
		{"b": 2401},
		{"c": 2500},
		{"d": 2304},
		{"e": 2401},
		{"f": 3364},
		{"g": 2809},
		{"h": 2304},
		{"i": 2401},
		{"j": 3025},
		{"k": 2916},
	}

	x, y, err := service.UnMapSlideCaptchaData(testV)
	if err != nil {
		t.Errorf("UnMapSlideCaptchaData error: <%s>\n", err.Error())
	}

	fmt.Println(x, y)
}
