package service_test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/internal/service"
	"uims/pkg/glog"
	"uims/pkg/thrift/client"
)

func TestThriftClientInvoke(t *testing.T) {
	glog.Init()
	p := struct {
		A int    `json:"a"`
		B string `json:"b"`
	}{
		A: 1,
		B: "test_get_api",
	}

	resp := service.GetThriftClientServer().ClientInvoke(1, "delete_resource_group", p)
	if !resp.OK() {
		t.Fatal(errors.New(resp.Err()))
	}
}

func TestThriftClientServer_VzhuoUserThriftClientInvoke(t *testing.T) {
	resp := service.GetThriftClientServer().ClientInvoke(
		3,
		"get_uims_wx_qr",
		map[string]interface{}{
			"params": map[string]interface{}{
				"items": nil,
			},
		},
	)
	assert.True(t, resp.OK())
}

func TestThriftClientServer_CASSThriftClientInvoke(t *testing.T) {
	resp := service.GetThriftClientServer().ClientInvoke(
		1,
		"test",
		map[string]interface{}{
			"params": map[string]interface{}{
				"items": nil,
			},
		},
	)
	assert.True(t, resp.OK())
}

func TestThriftClientServer_InvokeMP(t *testing.T) {
	type C struct {
		SceneId   float64 `json:"scene_id"`
		TicketUrl string  `json:"ticket_url"`
	}
	c := C{}
	r := service.Response{}
	service.ThriftClientServer{}.InvokeMP(service.Request{
		BRequest: thriftclient.BRequest{
			MethodName: "get_uims_wx_qr",
			Params: map[string]string{
				"test": "test",
			},
		},
	}, &r)
	t.Logf("%+v", r)
	t.Logf("%+v", r.Biz.BizContent)
	err := r.ParseContent(&c)
	assert.Nil(t, err)
	t.Logf("c: %+v", c)
	assert.True(t, r.OK())
	assert.NotEmpty(t, c.SceneId)
	assert.NotEmpty(t, c.TicketUrl)
	assert.Contains(t, c.TicketUrl, "http")
}
