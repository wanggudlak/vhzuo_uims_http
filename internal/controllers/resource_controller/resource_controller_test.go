package resource_controller_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"uims/app"
	"uims/boot"
	"uims/internal/controllers/resource_controller/requests"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	router = app.GetEngineRouter()
	os.Exit(m.Run())
}

func TestList(t *testing.T) {
	w := httptest.NewRecorder()

	reqForm := requests.ResourceListRequest{
		ClientId: 1,
	}
	requestString, _ := json.Marshal(reqForm)
	req, _ := http.NewRequest("GET", "/api/resource/list", strings.NewReader(string(requestString)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	t.Logf("%s \n", w.Body.String())
	res := struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Body    interface{} `json:"body"`
	}{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	t.Logf("%+v", res)
	assert.Nil(t, err)
	assert.Equal(t, 0, res.Code, res)
}
