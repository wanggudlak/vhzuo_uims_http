package user_controller_test

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uims/app"
	"uims/boot"
	"uims/internal/controllers/login_controller"
	responses2 "uims/internal/controllers/responses"
	requests2 "uims/internal/controllers/user_controller/requests"
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/randc"
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

func TestStore(t *testing.T) {
	request := requests2.UserStoreRequest{}
	_ = faker.FakeData(&request)
	request.Phone = randc.RandStringN(11)
	fmt.Printf("%+v", request)

	requestStr, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/users", strings.NewReader(string(requestStr)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	response := responses2.Parse(w.Body.Bytes())
	assert.Equal(t, responses2.CodeSuccess, response.Code)
	tool.Dump(response)
}

func TestList(t *testing.T) {
	var user model.User
	var err error
	err = db.Def().First(&user).Error
	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "api/users/list?phone=13517210601", strings.NewReader(string("")))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("token", login_controller.GenerateToken(user).Token)
	router.ServeHTTP(w, req)
	t.Logf("%+v", w)
	assert.Equal(t, http.StatusOK, w.Code)
	response := responses2.Parse(w.Body.Bytes())
	assert.Equal(t, responses2.CodeSuccess, response.Code)
	tool.Dump(response)
}
