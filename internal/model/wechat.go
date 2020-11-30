package model

type WeChat struct {
	ID     uint   `json:"id" gorm:"primary_key;comment:'主键'"`
	UUID   string `json:"uuid" gorm:"not null;type:varchar(32);column:uuid;comment:'唯一标识'"`
	AppId  string `json:"app_id" gorm:"default:'';comment:'客户端id, 对应表client'"`
	Secret string `json:"secret" gorm:"default:'';comment:'微信id, 对应表wechat'"`
	Desc   string `json:"desc" gorm:"default:'';comment:'描述'"`
	*CommonModel
}

func (WeChat) TableName() string {
	return "uims_wechat"
}
