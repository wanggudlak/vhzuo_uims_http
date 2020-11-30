package verifycode_controller_test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uims/app"
	"uims/boot"
	responses2 "uims/internal/controllers/responses"
	requests2 "uims/internal/controllers/verifycode_controller/requests"
	"uims/pkg/tool"
)

var router *gin.Engine
var w *httptest.ResponseRecorder

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	router = app.GetEngineRouter()
	w = httptest.NewRecorder()
	m.Run()
}

func TestGenerateMathCaptchaBase64(t *testing.T) {
	request := requests2.RequestOfGenerateMathCaptcha{}
	//_ = faker.FakeData(&request)
	fmt.Printf("%v\n", request)
	requestStr, _ := json.Marshal(request)
	req, _ := http.NewRequest("GET", "/api/captcha/math", strings.NewReader(string(requestStr)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	response := responses2.Parse(w.Body.Bytes())
	assert.Equal(t, responses2.CodeSuccess, response.Code)
	tool.Dump(response)
}
