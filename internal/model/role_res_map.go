package model

import "time"

type RoleResMap struct {
	ID           int       `gorm:"column:id;type:int(11) unsigned auto_increment;primary_key;;comment:'主键ID'" json:"id"`
	RoleID       int       `json:"role_id" gorm:"column:role_id;type:int(11) unsigned;not null;default:0;comment:'角色ID'"`
	ResGrpID     int       `gorm:"column:res_grp_id;type:int(11) unsigned;not null;default:0;comment:'角色关联的资源组ID'" json:"res_grp_id"`
	ClientID     int       `gorm:"column:client_id;type:int(11) unsigned;not null;default:0;comment:'客户端ID'" json:"client_id"`
	OrgId        int       `gorm:"column:org_id;type:int(11) unsigned;not null;default:0;comment:'组织ID，默认是0，即不区分组织'" json:"org_id"`
	StartValidAt time.Time `gorm:"column:start_valid_at;type:datetime;default:null;comment:'实际开始时间'" json:"start_valid_at"`
	IsDel        string    `gorm:"column:isdel;type:char(1);not null;default:'N';comment:'是否软删除，默认N：未软删除；Y：已软删除'" json:"isdel"`
	ForgetAt     time.Time `gorm:"column:forget_at;type:datetime;default:null;comment:'忘记时间'" json:"forget_at"`
	*CommonModel
}

// TableName sets the insert table name for this struct type
func (RoleResMap) TableName() string {
	return "uims_role_res_map"
}
