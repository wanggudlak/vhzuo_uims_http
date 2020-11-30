package model

type Org struct {
	ID          int    `gorm:"column:id;type:int(11) unsigned auto_increment;primary_key;comment:'组织ID'" json:"id"`
	ParentOrgID int    `gorm:"column:parent_org_id;type:int(11) unsigned;not null;default:0;comment:'直接父级组织ID'" json:"parent_org_id"`
	ClientID    uint   `gorm:"column:client_id;type:int(11) unsigned;not null;default:0;comment:'客户端ID'" json:"client_id"`
	ClientAppID string `gorm:"column:client_app_id;type:char(32);not null;default:'';comment:'客户端APPID'" json:"client_app_id"`
	OrgNameCN   string `gorm:"column:org_name_cn;type:varchar(255);not null;default:'';comment:'组织中文名'" json:"org_name_cn"`
	OrgNameEN   string `gorm:"column:org_name_en;type:varchar(255);not null;default:'';comment:'组织英文名'" json:"org_name_en"`
	//BusinessType   string `gorm:"column:business_type;type:varchar(32);not null;default:'';comment:'业务组织类型，cass：结算系统'" json:"business_type"`
	OrgCode        string `gorm:"column:org_code;type:varchar(255);not null;default:'';comment:'组织代码'" json:"org_code"`
	OrgLevel       int    `gorm:"column:org_level;type:tinyint(3) unsigned;not null;default:0;comment:'组织层级，0：顶级；1：第1级，以此类推'" json:"org_level"`
	OrgFullPinyin  string `gorm:"column:org_full_pinyin;type:varchar(255);not null;default:'';comment:'组织全拼音'" json:"org_full_pinyin"`
	OrgFirstPinyin string `gorm:"column:org_first_pinyin;type:varchar(255);not null;default:'';comment:'组织拼音搜字母'" json:"org_first_pinyin"`
	*CommonModel
}

// TableName sets the insert table name for this struct type
func (Org) TableName() string {
	return "uims_organization"
}
