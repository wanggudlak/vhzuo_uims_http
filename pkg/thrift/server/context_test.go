package thriftserver_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/boot"
	"uims/pkg/thrift/common"
	thriftserver "uims/pkg/thrift/server"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

func TestContextValidate(t *testing.T) {
	c := thriftserver.NewContext()
	c.Request.Body = `{"test": "test", "required_test": "a"}`
	type test struct {
		Test         string `json:"test" binding:"required"`
		RequiredTest string `json:"required_test" binding:"required" comment:"测试"`
	}
	var te test
	if err := c.ShouldBind(&te); err != nil {
		t.Logf("ShouldBind err: %+v", err)
		t.FailNow()
	}
	assert.Equal(t, te.Test, "test")
}

func TestContext_ParseRequest(t *testing.T) {
	c := thriftserver.NewContext()
	var err error
	// 无 method 参数
	err = c.ParseRequest(`{"test": "test"}`)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), common.INVALID_METHOD_NAME)

	// 有 method 参数
	err = c.ParseRequest(`{"method_name": "test","params": "{\"test\": \"test\", \"test_arr\": [\"a\"]}"}`)
	assert.Nil(t, err)
	type test struct {
		Method  string   `json:"method"`
		Test    string   `json:"test" binding:"required"`
		TestArr []string `json:"test_arr" binding:"required"`
	}
	var te test
	if err := c.ShouldBind(&te); err != nil {
		t.Logf("ShouldBind err: %+v", err)
		t.FailNow()
	}
	assert.Equal(t, te.Test, "test")
	assert.Equal(t, te.TestArr[0], "a")
}

func TestResponse_BadParams(t *testing.T) {
	c := thriftserver.NewContext()
	var err error
	if err = c.ParseRequest(`{"method_name": "test", "params": "{\"test\": \"a\"}"}`); err != nil {
		t.Error(err)
		t.FailNow()
	}
	type params struct {
		Test2 string `json:"test_2" binding:"required" comment:"测试字段"`
	}
	var p params
	if err = c.ShouldBind(&p); err != nil {
		c.Response.BadParams(err)
		assert.Equal(t, c.Response.Status, common.STATUS_SUCCESS)
		assert.Equal(t, c.Response.Message, common.CALL_SUCCESS_MSG)
		assert.Equal(t, c.Response.Data.Status, common.STATUS_PARAMS_FAILED)
		t.Logf("err: %+v", err)
		t.Logf("response: %+v", c.Response.ConvertThriftResp())
		assert.Equal(t, "测试字段为必填字段", c.Response.Data.Message)
		assert.Equal(t, common.STATUS_PARAMS_FAILED, c.Response.Data.Status)
		assert.Equal(t, "", c.Response.Data.Content)
	} else {
		t.Error("must has params parse error, but get nil")
		t.FailNow()
	}
}
func TestResponse_BadParams2(t *testing.T) {
	c := thriftserver.NewContext()
	var err error
	if err = c.ParseRequest(`{"method_name": "test", "params": ""}`); err != nil {
		t.Error(err)
		t.FailNow()
	}
	type params struct {
		Test2 string `json:"test_2" binding:"required" comment:"测试字段"`
	}
	var p params
	if err = c.ShouldBind(&p); err != nil {
		c.Response.BadParams(err)
		assert.Equal(t, c.Response.Status, common.STATUS_FAILED)
		assert.Equal(t, c.Response.Message, common.INVALID_PARAMS)
	} else {
		t.Error("must has params parse error, but get nil")
		t.FailNow()
	}
}
