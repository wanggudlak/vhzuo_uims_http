package model

import "time"

type UserWeChat struct {
	ID            uint       `json:"id" gorm:"primary_key;comment:'主键'"`
	UserId        uint       `json:"user_id" gorm:"INDEX;not null;comment:'用户id'"`
	WeChatId      uint       `json:"wechat_id" gorm:"column:wechat_id;comment:'wechat 配置表 id'"`
	Nickname      string     `json:"nickname" gorm:"column:nickname;type:varchar(40);default:''"`
	Sex           string     `json:"sex" gorm:"column:sex;type:char(1);default:'';comment:'M男 F女 空未知'"`
	Country       string     `json:"country" gorm:"column:country;type:varchar(20);default:''"`
	Avatar        string     `json:"avatar" gorm:"column:avatar;type:varchar(200);default:''"`
	Privilege     string     `json:"privilege" gorm:"column:privilege;type:varchar(150);default:''"`
	Province      string     `json:"province" gorm:"column:province;type:varchar(150);default:''"`
	City          string     `json:"city" gorm:"column:city;type:varchar(150);default:'';comment:'市'"`
	AccessToken   string     `json:"access_token" gorm:"column:access_token;type:varchar(150);default:''"`
	RefreshToken  string     `json:"refresh_token" gorm:"column:refresh_token;type:varchar(150);default:''"`
	WeChatOpenId  string     `json:"wechat_open_id" gorm:"column:wechat_open_id;UNIQUE;comment:'微信 openId';type:varchar(40)"`
	WeChatUnionId string     `json:"wechat_union_id" gorm:"column:wechat_union_id;INDEX;comment:'微信 unionId';type:varchar(40)"`
	DeletedAt     *time.Time `sql:"index"`
	*CommonModel
}

func (UserWeChat) TableName() string {
	return "uims_user_wechat"
}
