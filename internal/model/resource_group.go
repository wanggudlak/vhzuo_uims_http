package model

import (
	"database/sql/driver"
	"encoding/json"
)

type ResourceGroup struct {
	ID           uint            `gorm:"column:id;primary_key;comment:'资源组ID'" json:"id"`
	ResGroupCode string          `gorm:"column:res_group_code;type:varchar(32);not null;default:'';index:res_group_code;comment:'权限资源策略组编码'" json:"res_group_code"`
	ResGroupEn   string          `gorm:"column:res_group_en;type:varchar(64);not null;default:'';comment:'权限策略组英文名称'" json:"res_group_en"`
	ResGroupCn   string          `gorm:"column:res_group_cn;type:varchar(64);not null;default:'';comment:'权限策略组中文名称'" json:"res_group_cn"`
	ResGroupType string          `gorm:"column:res_group_type;type:char(10);not null;default:'DEFAULT';comment:'权限策略组类型：DEFAULT-默认策略组；SELF-自定义配置的策略组'" json:"res_group_type"`
	ResOfCurr    *ResourceOfCurr `gorm:"column:res_of_curr;type:json;comment:'属于当前策略组的资源id list'" json:"res_of_curr"`
	ClientId     uint            `gorm:"column:client_id;type:int(11) unsigned;not null;default:0;index:client_id;comment:'客户端或业务系统ID，默认是0，即不区分客户端业务系统，属于跨业务系统通用类型策略组'" json:"client_id"`
	OrgId        uint            `gorm:"column:org_id;type:int(11) unsigned;not null;default:0;index:org_id;comment:'组织id，默认是0，标识不区分组织，即是跨组织型的策略组'" json:"org_id"`
	IsDel        string          `gorm:"column:isdel;type:char(1);not null;default:'N';comment:'是否软删除，默认N：未软删除；Y：已软删除'" json:"isdel"`
	*CommonModel
}

type ResourceOfCurr struct {
	ResourceIDs []int `gorm:"column:resource_ids" json:"resource_ids"`
}

// TableName sets the insert table name for this struct type
func (r ResourceGroup) TableName() string {
	return "uims_res_group"
}

// json类型必须实现Value和Scan方法
func (r *ResourceOfCurr) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan 实现方法
func (r *ResourceOfCurr) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &r)
}
