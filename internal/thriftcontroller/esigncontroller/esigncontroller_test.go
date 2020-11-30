package esigncontroller_test

import (
	"fmt"
	genid "github.com/srlemon/gen-id"
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/boot"
	"uims/internal/model"
	"uims/internal/thriftcontroller/esigncontroller"
	"uims/pkg/db"
	"uims/pkg/thrift/common"
	thriftserver "uims/pkg/thrift/server"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

func TestNotifyESignExistUser(t *testing.T) {
	var err error
	c := thriftserver.NewContext()
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: esigncontroller.NotifyESignReq{
			Name:                        "詹光",
			IdCard:                      "420222199212041057",
			Phone:                       "13517210601",
			IdentityCardPersonImgBase64: "",
			IdentityCardEmblemImgBase64: "",
		},
	}.JSON())
	esigncontroller.NotifyESign(c)
	assert.Nil(t, err)
	t.Logf("response: %+v", c.Response)
	t.Logf("response data: %+v", c.Response.Data)

	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	if resp, ok := c.Response.Data.Content.(esigncontroller.NotifyESignResp); !ok {
		t.Log(resp)
		t.Errorf("%+v type not is %s", c.Response.Data.Content, "esigncontroller.NotifyESignResp")
		t.FailNow()
	} else {
	}
	t.Logf("response thrift: %+v", c.Response.ConvertThriftResp())
}

func TestNotifyESignNoExistUser(t *testing.T) {
	var err error
	c := thriftserver.NewContext()
	name := genid.NewGeneratorData().Name
	idCard := genid.NewGeneratorData().IDCard
	phone := genid.NewGeneratorData().PhoneNum
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: esigncontroller.NotifyESignReq{
			Name:                        name,
			IdCard:                      idCard,
			Phone:                       phone,
			IdentityCardPersonImgBase64: "",
			IdentityCardEmblemImgBase64: "",
		},
	}.JSON())
	esigncontroller.NotifyESign(c)
	assert.Nil(t, err)
	t.Logf("response: %+v", c.Response)
	t.Logf("response data: %+v", c.Response.Data)

	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	if resp, ok := c.Response.Data.Content.(esigncontroller.NotifyESignResp); !ok {
		t.Log(resp)
		t.Errorf("%+v type not is %s", c.Response.Data.Content, "esigncontroller.NotifyESignResp")
		t.FailNow()
	} else {
		// 执行成功
		var user model.User
		err = db.Def().Where(&model.User{
			Phone:          &phone,
			IdentityCardNo: &idCard,
		}).First(&user).Error
		assert.Nil(t, err)
		assert.NotEmpty(t, user.ID)
	}
}

func TestNotifyESignExistPhone(t *testing.T) {
	var err error
	c := thriftserver.NewContext()
	name := genid.NewGeneratorData().Name
	idCard := genid.NewGeneratorData().IDCard
	phone := genid.NewGeneratorData().PhoneNum
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: esigncontroller.NotifyESignReq{
			Name:                        name,
			IdCard:                      idCard,
			Phone:                       phone,
			IdentityCardPersonImgBase64: "",
			IdentityCardEmblemImgBase64: "",
		},
	}.JSON())
	esigncontroller.NotifyESign(c)
	assert.Nil(t, err)
	t.Logf("response: %+v", c.Response)
	t.Logf("response data: %+v", c.Response.Data)

	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	if resp, ok := c.Response.Data.Content.(esigncontroller.NotifyESignResp); !ok {
		t.Log(resp)
		t.Errorf("%+v type not is %s", c.Response.Data.Content, "esigncontroller.NotifyESignResp")
		t.FailNow()
	} else {
		// 执行成功
		var user model.User
		err = db.Def().Where(&model.User{
			Phone:          &phone,
			IdentityCardNo: &idCard,
		}).First(&user).Error
		assert.Nil(t, err)
		assert.NotEmpty(t, user.ID)
	}

	// 传入手机号不同, 身份证号一样的数据
	newPhone := genid.NewGeneratorData().PhoneNum
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: esigncontroller.NotifyESignReq{
			Name:                        name,
			IdCard:                      idCard,
			Phone:                       newPhone,
			IdentityCardPersonImgBase64: "",
			IdentityCardEmblemImgBase64: "",
		},
	}.JSON())
	esigncontroller.NotifyESign(c)
	assert.Nil(t, err)
	t.Logf("response: %+v", c.Response)
	t.Logf("response data: %+v", c.Response.Data)

	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	assert.Equal(t, common.STATUS_FAILED, c.Response.Data.Status)
	assert.Equal(t, fmt.Sprintf("传入的手机号 %s 和系统中已存在的手机号 %s 不一致", newPhone, phone), c.Response.Data.Message)

	// 传入身份证号不同, 手机号一样的数据
	newIdCard := genid.NewGeneratorData().IDCard
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: esigncontroller.NotifyESignReq{
			Name:                        name,
			IdCard:                      newIdCard,
			Phone:                       phone,
			IdentityCardPersonImgBase64: "",
			IdentityCardEmblemImgBase64: "",
		},
	}.JSON())
	esigncontroller.NotifyESign(c)
	assert.Nil(t, err)
	t.Logf("response: %+v", c.Response)
	t.Logf("response data: %+v", c.Response.Data)

	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	assert.Equal(t, common.STATUS_FAILED, c.Response.Data.Status)
	assert.Equal(t, fmt.Sprintf("传入的身份证号 %s 和系统中已存在的身份证号 %s 不一致", newIdCard, idCard), c.Response.Data.Message)
}

func TestNotifyESignNoPhone(t *testing.T) {
	var err error
	c := thriftserver.NewContext()
	err = c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: esigncontroller.NotifyESignReq{
			Name:                        genid.NewGeneratorData().Name,
			IdCard:                      genid.NewGeneratorData().IDCard,
			Phone:                       "",
			IdentityCardPersonImgBase64: "",
			IdentityCardEmblemImgBase64: "",
		},
	}.JSON())
	esigncontroller.NotifyESign(c)
	assert.Nil(t, err)
	t.Logf("response: %+v", c.Response)
	t.Logf("response data: %+v", c.Response.Data)

	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	if resp, ok := c.Response.Data.Content.(esigncontroller.NotifyESignResp); !ok {
		t.Log(resp)
		t.Errorf("%+v type not is %s", c.Response.Data.Content, "esigncontroller.NotifyESignResp")
		t.FailNow()
	} else {
	}
	t.Logf("response thrift: %+v", c.Response.ConvertThriftResp())
}
