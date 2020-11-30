package model

import "time"

type UserInfo struct {
	ID                    int        `gorm:"column:id;type:int(11) unsigned auto_increment;primary_key;comment:'主键ID'" json:"id"`
	UserID                int        `gorm:"column:user_id;type:bigint(11) unsigned;not null;comment:'用户ID'" json:"user_id"`
	IsIdentify            string     `gorm:"column:is_identify;type:char(1);default:'';comment:'是否实名认证，默认N：没有；Y：已经实名认证'" json:"is_identify"`
	UserCode              string     `gorm:"column:user_code;type:char(19);not null;default:'';comment:'用户编码，全数字最多11位，不同的组织下可以重复'" json:"user_code"`
	UserType              string     `gorm:"column:user_type;type:char(32);not null;default:'';comment:'用户类型，取client_type字段的值，VDK：微桌，CASS：结算系统'" json:"user_type"`
	UserBussType          string     `gorm:"column:user_buss_type;type:varchar(32);not null;default:'normal';comment:'用户业务类型，normal：普通用户，business：商户，settle_company：结算公司，back：后台用户，uims：uims用户'" json:"user_buss_type"`
	NameEn                string     `gorm:"column:name_en;type:varchar(64);not null;default:'';comment:'用户英文姓名'" json:"name_en"`
	NameCn                string     `gorm:"column:name_cn;type:varchar(64);not null;default:'';comment:'用户中文名'" json:"name_cn"`
	NameCnAlias           string     `gorm:"column:name_cn_alias;type:varchar(16);not null;default:'';comment:'用户别名'" json:"name_cn_alias"`
	NameAbbrPy            string     `gorm:"column:name_abbr_py;type:varchar(16);default:'';comment:'用户姓名拼音首字母'" json:"name_abbr_py"`
	NameFullPy            string     `gorm:"column:name_full_py;type:varchar(255);comment:'用户姓名全拼音，英文空格分隔'" json:"name_full_py"`
	IdentityCardNo        *string    `gorm:"column:identity_card_no;type:varchar(20);unique;comment:'用户的身份证号'" json:"identity_card_no"`
	NaCode                string     `gorm:"column:na_code;type:char(5);not null;default:'';comment:'国家代码，中国：+86'" json:"na_code"`
	Phone                 string     `gorm:"column:phone;type:char(12);not null;default:'';comment:'手机号'" json:"phone"`
	LandlinePhone         string     `gorm:"column:landline_phone;type:varchar(16);not null;default:'';comment:'座机号'" json:"landline_phone"`
	Sex                   string     `gorm:"column:sex;type:char(1);not null;default:'';comment:'性别，M：男；F：女'" json:"sex"`
	Birthday              *time.Time `gorm:"column:birthday;type:date;default:null;comment:'出生年月日'" json:"birthday" format:"2006-01-02 15:04:05"`
	Nickname              string     `gorm:"column:nickname;type:varchar(32);not null;default:'';comment:'昵称'" json:"nickname"`
	TaxerType             string     `gorm:"column:taxer_type;type:char(1);not null;default:'';comment:'纳税人类型，A：一般纳税人'" json:"taxer_type"`
	TaxerNo               string     `gorm:"column:taxer_no;type:varchar(16);not null;default:'';comment:'纳税人识别号'" json:"taxer_no"`
	HeaderImgURL          string     `gorm:"column:header_img_url;type:varchar(255);not null;default:'';comment:'用户头像图片相对地址'" json:"header_img_url"`
	IdentityCardPersonImg string     `gorm:"column:identity_card_person_img;type:varchar(255);not null;default:'';comment:'用户身份证人像面图片相对地址'" json:"identity_card_person_img"`
	IdentityCardEmblemImg string     `gorm:"column:identity_card_emblem_img;type:varchar(255);not null;default:'';comment:'用户身份证人像面图片相对地址'" json:"identity_card_emblem_img"`
	Isdel                 string     `gorm:"column:isdel;type:char(1);not null;default:'N';comment:'是否软删除，默认N：未软删除；Y：已软删除'" json:"-"`
	*CommonModel
}

// TableName sets the insert table name for this struct type
func (UserInfo) TableName() string {
	return "uims_user_info"
}
