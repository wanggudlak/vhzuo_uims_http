package model

type UserRole struct {
	ID               int    `gorm:"column:id;type:int(11) unsigned auto_increment;primary_key;;comment:'主键ID'" json:"id"`
	UserID           int    `gorm:"column:user_id;type:bigint(11) unsigned;not null;default:0;comment:'用户ID'" json:"user_id"`
	RoleID           int    `gorm:"column:role_id;type:int(11) unsigned;not null;default:0;comment:'角色ID'" json:"role_id"`
	ClientID         int    `gorm:"column:client_id;type:int(11) unsigned;not null;default:0;comment:'客户端ID'" json:"client_id"`
	UserRelationType string `gorm:"column:user_relation_type;type:varchar(12);not null;default:'user';comment:'用户关联类型，user：用户，org：组织'" json:"user_relation_type"`
	*CommonModel
}

// TableName sets the insert table name for this struct type
func (UserRole) TableName() string {
	return "uims_user_role"
}
