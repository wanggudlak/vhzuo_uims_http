package wechatqrcode

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"uims/pkg/gjwt"
	"uims/pkg/gredis"
	"uims/pkg/tool"
)

const TTL = 5 * time.Minute

type BindPhoneClaims struct {
	gjwt.Jwt
	WeChatInfo WeChatInfo `json:"wechat_info"`
}

// {'province': '湖北', 'city': '武汉', 'groupid': 0, 'sex': 1, 'country': '中国', 'tagid_list': [], 'language': 'zh_CN', 'remark': '', 'headimgurl': 'http://thirdwx.qlogo.cn/mmopen/GaUReTcGVTvQCrfMjwhPgicGvAWicxRkyQU8dCR2EXUNPCqKOUS8CrprePQ4Q0qVkxwicLuXicFKP0va9FMYgBmF7icGHic49X4usq/132', 'subscribe_time': 1592965663, 'nickname': '艾艾艾', 'qr_scene_str': '', 'subscribe': 1, 'subscribe_scene': 'ADD_SCENE_QR_CODE', 'qr_scene': 61347905, 'openid': 'ojyIo55IFpDmsYJPJP-3f7P9ZuTQ'}
type WeChatInfo struct {
	City       string `json:"city"`
	Province   string `json:"province"`
	Sex        int    `json:"sex"`
	Country    string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
	Nickname   string `json:"nickname"`
	OpenID     string `json:"openid"`
}

// 设置场景
func SetScene(sceneId int) error {
	key := GetRedisKey(strconv.Itoa(sceneId))
	err := gredis.Def().HSet(key, map[string]interface{}{
		"time": time.Now().Unix(),
	}).Err()
	if err != nil {
		return err
	}
	return gredis.Def().Expire(key, TTL).Err()
}

// 设置已经扫码了
func SetScanOK(sceneId int) error {
	key := GetRedisKey(strconv.Itoa(sceneId))
	err := gredis.Def().HSet(key, map[string]interface{}{
		"scan": true,
	}).Err()
	if err != nil {
		return errors.Wrap(err, "hset 失败")
	}
	return nil
}

// 是否已经扫码
func IsScan(sceneId int) bool {
	key := GetRedisKey(strconv.Itoa(sceneId))
	v := gredis.Def().HGet(key, "scan").Val()
	return v == "1"
}

// 设置场景授权成功
func SetSceneAuthOK(sceneId int, userId uint) error {
	key := GetRedisKey(strconv.Itoa(sceneId))
	if gredis.Def().Exists(key).Val() == 0 {
		return errors.New(fmt.Sprintf("scene_id key %s 不存在", key))
	}
	err := gredis.Def().HSet(key, map[string]interface{}{
		"user_id":  userId,
		"login_ok": true,
	}).Err()
	if err != nil {
		return errors.Wrap(err, "hset 失败")
	}
	return nil
}

// 设置需要绑定微信
func SetNeedBindPhone(sceneId int, weChatInfo WeChatInfo) error {
	key := GetRedisKey(strconv.Itoa(sceneId))
	return gredis.Def().HSet(key, "wechat_info", tool.JSONString(weChatInfo), "need_bind_phone", true).Err()
}

// 是否需要绑定
func NeedBindPhone(sceneId int) bool {
	key := GetRedisKey(strconv.Itoa(sceneId))
	v := gredis.Def().HGet(key, "need_bind_phone").Val()
	return v == "1"
}

// 获取绑定保存的微信信息
func GetNeedBindPhoneWeChatInfo(sceneId int) (*WeChatInfo, error) {
	key := GetRedisKey(strconv.Itoa(sceneId))
	v := gredis.Def().HGet(key, "wechat_info").Val()
	w := WeChatInfo{}
	err := json.Unmarshal([]byte(v), &w)
	if err != nil {
		return nil, errors.Wrap(err, "获取微信信息失败")
	}
	return &w, nil
}

// 场景id是否存在
func Exists(sceneID int) bool {
	key := GetRedisKey(strconv.Itoa(sceneID))
	v := gredis.Def().Exists(key).Val()
	return v == 1
}

// 获取场景值
func GetSceneUserID(sceneId int) (int, error) {
	key := GetRedisKey(strconv.Itoa(sceneId))
	i, err := strconv.Atoi(gredis.Def().HGet(key, "user_id").Val())
	if err != nil {
		return 0, errors.Wrap(err, "查询场景值失败")
	}
	return i, nil
}

func GetRedisKey(key string) string {
	return "uims:wechat:qrcode:login:sceneid:" + key
}
