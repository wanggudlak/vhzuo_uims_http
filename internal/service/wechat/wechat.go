package wechat

import (
	"uims/internal/model"
	"uims/pkg/db"
	"uims/pkg/wechatserver"
)

func GetConfig(clientId int) (*wechatserver.Config, error) {
	var err error
	var config = &wechatserver.Config{}
	c := &model.ClientWeChat{}
	err = db.Def().Select([]string{"id", "wechat_id"}).
		Where("client_id = ?", clientId).
		First(&c).Error
	if err != nil {
		return nil, err
	}
	w := &model.WeChat{}
	err = db.Def().Select([]string{"id", "app_id", "secret"}).
		Where("id = ?", c.WeChatId).
		First(&w).Error
	if err != nil {
		return nil, err
	}
	config.AppId = w.AppId
	config.Secret = w.Secret
	config.WeChatId = c.WeChatId
	config.ClientId = c.ClientId
	return config, nil
}

func GetConfigByWeChatId(weChatId uint) (*wechatserver.Config, error) {
	var config = &wechatserver.Config{}
	w := &model.WeChat{}
	err := db.Def().Select([]string{"id", "app_id", "secret"}).
		Where("id = ?", weChatId).
		First(&w).Error
	if err != nil {
		return nil, err
	}
	config.AppId = w.AppId
	config.Secret = w.Secret
	config.WeChatId = weChatId
	config.ClientId = 0
	return config, nil
}

func WeChatUUID2WeChatId(weChatUUID string) (uint, error) {
	var w model.WeChat
	err := db.Def().Select("id").Where("uuid = ?", weChatUUID).First(&w).Error
	if err != nil {
		return 0, err
	}
	return w.ID, nil
}
