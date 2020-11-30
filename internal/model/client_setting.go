package model

import (
	"database/sql/driver"
	"encoding/json"
)

type ClientSetting struct {
	ID           uint   `gorm:"column:id;type:int(11) unsigned auto_increment;primary_key;comment:'主键ID'" json:"id"`
	ClientID     uint   `gorm:"column:client_id;type:int(11) unsigned;not null;default:0;comment:'客户端ID'" json:"client_id"`
	Type         string `gorm:"column:type;type:char(3);not null;default:'';comment:'类型：LGN-用于登录的设置；REG-用于注册的设置；'" json:"type"`
	BusChannelID string `gorm:"column:bus_channel_id;type:char(3);not null;default:'';comment:'频道ID，对于登录业务，频道ID为100；注册业务频道ID为200'" json:"bus_channel_id"`
	PageID       string `gorm:"column:page_id;type:char(3);not null;default:'';comment:'页面ID，对于登录业务的登录页面ID为101；注册业务的注册页面ID为201'" json:"page_id"`
	SpmFullCode  string `gorm:"column:spm_full_code;type:char(32);not null;default:'';comment:'spm编码，由以下组成：client_spm1_code.client_spm2_code.频道ID.页面ID'" json:"spm_full_code"`
	// comment:'"表单域属性数据{"a": [{"attr_id": "account", "attr_cn": "账号"},{"attr_id": "passwd", "attr_cn": "密码"}, {"attr_id": "sms_code", "attr_cn": "验证码", "type": "phone_sms"}]}
	FormFields *FieldsMap `gorm:"column:form_fields;type:json;default:null;'" json:"form_fields"`
	//PageTemplateFile string     `gorm:"column:page_template_file;type:json;default:null;comment:'页面路径'" json:"page_template_file"`
	PageTemplateFile *PageTemplateFile `gorm:"column:page_template_file;type:json;default:null;comment:'页面路径'" json:"page_template_file"`
	Isdel            string            `gorm:"column:isdel;type:char(1);not null;default:'N';comment:'是否软删除，默认N：未删除；Y：已软删除'" json:"isdel"`
	*CommonModel
}

// TableName sets the insert table name for this struct type
func (ClientSetting) TableName() string {
	return "uims_client_settings"
}

func (c *ClientSetting) TemplateFile() string {
	return c.PageTemplateFile.A
}

type FieldsMap struct {
	Src   map[string]interface{}
	Valid bool
}

type PageTemplateFile struct {
	A string `json:"a"`
	B string `json:"b"`
}

//func NewFieldsMap(src map[string]interface{}) *FieldsMap {
//	if src != nil {
//		return &FieldsMap{
//			Src:   src,
//			Valid: true,
//		}
//	} else {
//		return &FieldsMap{
//			Src:   make(map[string]interface{}),
//			Valid: true,
//		}
//	}
//}

func (ls *FieldsMap) Scan(value interface{}) error {
	if value == nil {
		ls.Src, ls.Valid = make(map[string]interface{}), false
		return nil
	}
	t := make(map[string]interface{})
	if e := json.Unmarshal(value.([]byte), &t); e != nil {
		return e
	}
	ls.Valid = true
	ls.Src = t
	return nil
}

func (ls *FieldsMap) Value() (driver.Value, error) {
	if ls == nil {
		return nil, nil
	}
	if !ls.Valid {
		return nil, nil
	}

	b, e := json.Marshal(ls.Src)
	return b, e
}

// Value 实现方法
func (p *PageTemplateFile) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *PageTemplateFile) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), p)
}
