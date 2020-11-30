package model

type UserOrg struct {
	ID       int `gorm:"column:id;type:int(11) unsigned auto_increment;primary_key;comment:'主键ID'" json:"id"`
	UserID   int `gorm:"column:user_id;type:bigint(11) unsigned;default:0;not null;comment:'用户ID'" json:"user_id"`
	ClientID int `gorm:"column:client_id;type:int(11) unsigned;not null;default:0;comment:'客户端ID'" json:"client_id"`
	OrgID    int `gorm:"column:org_id;type:int(11);not null;default:0;comment:'客户端组织ID'" json:"org_id"`
	*CommonModel
}

// TableName sets the insert table name for this struct type
func (UserOrg) TableName() string {
	return "uims_user_org"
}
