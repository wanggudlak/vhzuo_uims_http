package requests

import "uims/internal/model"

type ResourceCreateRequest struct {
	ClientId        uint               `json:"client_id" form:"client_id" binding:"required"`
	OrgId           uint               `json:"org_id" form:"org_id" binding:"number"`
	ResType         string             `json:"res_type" form:"client_id" binding:"required"`
	ResSubType      string             `json:"res_sub_type" form:"res_sub_type" binding:"required"`
	ResNameEn       string             `json:"res_name_en" form:"res_name_en" binding:"required"`
	ResNameCn       string             `json:"res_name_cn" form:"res_name_cn" binding:"required"`
	ResEndpRoute    string             `json:"res_endp_route" form:"res_endp_route" binding:"-"`
	ResDataLocation model.LocationData `json:"res_data_location" form:"res_data_location" binding:"-"`
}
