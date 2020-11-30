package model

import (
	"database/sql/driver"
	"encoding/json"
)

type Resource struct {
	ID              uint          `gorm:"column:id;primary_key;comment:'资源ID'" json:"id"`
	ClientId        uint          `gorm:"column:client_id;type:int(11);not null;default:0;index:client_id;comment:'客户端业务系统ID'" json:"client_id"`
	OrgId           uint          `gorm:"column:org_id;type:int(11);not null;default:0;index:org_id;comment:'客户端组织ID'" json:"org_id"`
	ResCode         string        `gorm:"column:res_code;type:varchar(32);not null;default:'';index:res_code;comment:'资源编码'" json:"res_code"`
	ResFrontCode    string        `gorm:"column:res_front_code;type:varchar(32);not null;default:'';index:res_front_code;comment:'和前端约定的资源编码'" json:"res_front_code"`
	ResType         string        `gorm:"column:res_type;type:char(1);not null;default:'';comment:'资源类型，A：逻辑资源；B：实体资源'" json:"res_type"`
	ResSubType      string        `gorm:"column:res_sub_type;type:char(3);not null;default:'';comment:'资源子类型，AP：页面；AC：菜单；AM：按钮；AD：数据资源'" json:"res_sub_type"`
	ResNameEn       string        `gorm:"column:res_name_en;type:varchar(64);not null;default:'';comment:'资源的英文名称'" json:"res_name_en"`
	ResNameCn       string        `gorm:"column:res_name_cn;type:varchar(64);not null;default:'';comment:'资源的中文名称'" json:"res_name_cn"`
	ResEndpRoute    string        `gorm:"column:res_endp_route;type:varchar(128);not null;default:'';comment:'资源的后端路由URI'" json:"res_endp_route"`
	ResDataLocation *LocationData `gorm:"column:res_data_location;type:json;comment:'资源所在的位置，主要用于数据权限，json存储，包含以下属性：客户端id、数据库名、表名、行记录属性名、行记录属性值'" json:"res_data_location"`
	IsDel           string        `gorm:"column:isdel;type:char(1);not null;default:'N';comment:'是否软删除，默认N：未软删除；Y：已软删除'" json:"isdel"`
	*CommonModel
}

func (Resource) TableName() string {
	return "uims_access_resource"
}

type LocationData struct {
	DataBase string `json:"database"`
	Table    string `json:"table"`
	Status   string `json:"status"`
}

// json类型必须实现Value和Scan方法
func (r LocationData) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan 实现方法
func (r *LocationData) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &r)
}
