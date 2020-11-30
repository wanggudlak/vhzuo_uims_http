package model

type Role struct {
	ID         int    `gorm:"column:id;type:int(11) unsigned auto_increment;primary_key;comment:'角色ID'" json:"id"`
	ClientID   int    `gorm:"column:client_id;type:int(11) unsigned;not null;default:0;comment:'客户端ID'" json:"client_id"`
	OrgID      int    `gorm:"column:org_id;type:int(11) unsigned;not null;default:0;comment:'所属的组织ID'" json:"org_id"`
	RoleType   string `gorm:"column:role_type;type:char(1);not null;default:'';comment:'角色类型：A：UIMS的角色；F：通过页面增加的角色；C：结算系统；V：微桌系统'" json:"role_type"` // 角色类型：A：UIMS的角色；F：通过页面增加的角色
	RoleCode   string `gorm:"column:role_code;type:varchar(32);not null;default:'';comment:'资源编码，UIMS系统的角色固定用UIMS.SUPERADMIN.001'" json:"role_code"` //资源编码，UIMS系统的角色固定用UIMS.SUPERADMIN.001
	RoleNameEN string `gorm:"column:role_name_en;type:varchar(64);not null;default:'';comment:'资源英文名称'" json:"role_name_en"`
	RoleNameCN string `gorm:"column:role_name_cn;type:varchar(64);not null;default:'';comment:'资源中文名称'" json:"role_name_cn"`
	IsDel      string `gorm:"column:isdel;type:char(1);not null;default:'N';comment:'是否软删除，默认N：未软删除；Y：已软删除'" json:"isdel"` //是否软删除，默认N：未软删除；Y：已软删除
	*CommonModel
}

func (Role) TableName() string {
	return "uims_role"
}
