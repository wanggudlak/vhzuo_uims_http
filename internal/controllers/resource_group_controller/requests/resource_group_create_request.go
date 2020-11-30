package requests

import "uims/internal/model"

type ResourceGroupCreateRequest struct {
	ClientId     uint                  `json:"client_id" from:"client_id" binding:"required"`
	OrgId        uint                  `json:"org_id" from:"org_id" binding:"number"`
	ResGroupCn   string                `json:"res_group_cn" from:"res_group_cn" biding:"required"`
	ResGroupEn   string                `json:"res_group_en" from:"res_group_en" biding:"required"`
	ResGroupCode string                `json:"res_group_code" from:"res_group_code" biding:"required"`
	ResGroupType string                `json:"res_group_type" from:"res_group_type" biding:"required"`
	ResOfCurr    *model.ResourceOfCurr `json:"res_of_curr" from:"res_of_curr" binding:"required"`
}
