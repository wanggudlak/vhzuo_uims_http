package wechatcontroller_test

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"uims/boot"
	"uims/internal/service/wechatqrcode"
	"uims/internal/thriftcontroller/vzhuo/wechatcontroller"
	"uims/pkg/randc"
	"uims/pkg/thrift/common"
	thriftserver "uims/pkg/thrift/server"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

func TestFollowEvent(t *testing.T) {
	c := thriftserver.NewContext()
	sceneID := rand.Intn(999999999)
	err := c.ParseRequest(thriftserver.BaseRequest{
		Method: "test",
		Params: wechatcontroller.FollowEventReq{
			SceneId: sceneID,
			WXInfo: wechatqrcode.WeChatInfo{
				City:       "武汉",
				Province:   "湖北",
				Sex:        1,
				Country:    "中国",
				HeadImgURL: "http://thirdwx.qlogo.cn/mmopen/GaUReTcGVTvQCrfMjwhPgicGvAWicxRkyQU8dCR2EXUNPCqKOUS8CrprePQ4Q0qVkxwicLuXicFKP0va9FMYgBmF7icGHic49X4usq/132",
				Nickname:   "艾艾艾",
				OpenID:     randc.RandStringN(28),
			},
		},
	}.JSON())
	assert.Nil(t, err)
	wechatcontroller.FollowEvent(c)
	assert.Equal(t, common.STATUS_SUCCESS, c.Response.Status)
	assert.Equal(t, common.CALL_SUCCESS_MSG, c.Response.Message)
	assert.Equal(t, common.STATUS_FAILED, c.Response.Data.Status)
}
