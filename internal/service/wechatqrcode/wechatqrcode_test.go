package wechatqrcode_test

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"uims/boot"
	"uims/internal/service/wechatqrcode"
	"uims/pkg/randc"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

func TestSetSceneAuthOK(t *testing.T) {
	var err error
	// sceneID := rand.Intn(999999999)
	sceneID := 1595312118785
	err = wechatqrcode.SetScene(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetScanOK(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetSceneAuthOK(sceneID, 1)
	assert.Nil(t, err)
}

func TestGetNeedBindPhoneWeChatInfo(t *testing.T) {
	sceneID := rand.Intn(999999999)
	err := wechatqrcode.SetScene(sceneID)
	assert.Nil(t, err)
	w1 := wechatqrcode.WeChatInfo{
		City:       "武汉",
		Province:   "湖北",
		Sex:        1,
		Country:    "中国",
		HeadImgURL: "",
		Nickname:   "测试",
		OpenID:     randc.RandStringN(28),
	}
	err = wechatqrcode.SetScanOK(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetNeedBindPhone(sceneID, w1)
	assert.Nil(t, err)
	w2, err := wechatqrcode.GetNeedBindPhoneWeChatInfo(sceneID)
	assert.Nil(t, err)
	assert.Equal(t, w1.OpenID, w2.OpenID)
}

func TestSetNeedBindPhone(t *testing.T) {
	sceneID := 1595311927668
	err := wechatqrcode.SetScene(sceneID)
	assert.Nil(t, err)
	w1 := wechatqrcode.WeChatInfo{
		City:       "武汉",
		Province:   "湖北",
		Sex:        1,
		Country:    "中国",
		HeadImgURL: "",
		Nickname:   "测试",
		OpenID:     randc.RandStringN(28),
	}
	err = wechatqrcode.SetScanOK(sceneID)
	assert.Nil(t, err)
	err = wechatqrcode.SetNeedBindPhone(sceneID, w1)
	assert.Nil(t, err)
	w2, err := wechatqrcode.GetNeedBindPhoneWeChatInfo(sceneID)
	assert.Nil(t, err)
	assert.Equal(t, w1.OpenID, w2.OpenID)
}

func TestExists(t *testing.T) {
	sceneID := 123456789
	assert.False(t, wechatqrcode.Exists(sceneID))
	err := wechatqrcode.SetScene(sceneID)
	assert.Nil(t, err)
	assert.True(t, wechatqrcode.Exists(sceneID))
}
