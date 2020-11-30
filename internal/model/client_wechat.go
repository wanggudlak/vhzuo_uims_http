package model

type ClientWeChat struct {
	ID       uint `json:"id" gorm:"primary_key;comment:'主键'"`
	ClientId uint `json:"client_id" gorm:"index:index_client_id;not null;comment:'客户端id, 对应表client'"`
	WeChatId uint `json:"wechat_id" gorm:"column:wechat_id;index:index_wechat_id;not null;comment:'微信id, 对应表wechat'"`
	*CommonModel
}

func (ClientWeChat) TableName() string {
	return "uims_client_wechat"
}
