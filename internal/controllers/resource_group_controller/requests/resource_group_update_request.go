package requests

import (
	"uims/internal/model"
)

type ResourceGroupUpdateRequest struct {
	ID           uint                  `json:"id" form:"id" binding:"required"`
	ResGroupEn   string                `json:"res_group_en" form:"res_group_en" biding:"required"`
	ResGroupCn   string                `json:"res_group_cn" form:"res_group_cn" biding:"required"`
	ResGroupType string                `json:"res_group_type" form:"res_group_type" biding:"required"`
	ResOfCurr    *model.ResourceOfCurr `json:"res_of_curr" form:"res_of_curr" binding:"required"`
}
