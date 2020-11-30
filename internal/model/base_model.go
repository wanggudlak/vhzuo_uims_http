package model

import "time"

type CommonModel struct {
	CreatedAt time.Time `json:"created_at" format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `json:"updated_at" format:"2006-01-02 15:04:05"`
}

// 初始化创建时间和更新时间字段
//func (cm *CommonModel) InitTime() *CommonModel {
//	cm.CreatedAt = time.Now()
//	cm.UpdatedAt = time.Now()
//	return cm
//}
