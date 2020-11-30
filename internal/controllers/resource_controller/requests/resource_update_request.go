package requests

import "uims/internal/model"

type ResourceUpdateRequest struct {
	ResourceId      int                `json:"id" form:"id" binding:"required"`
	ResType         string             `json:"res_type" form:"client_id" binding:"required"`
	ResSubType      string             `json:"res_sub_type" form:"res_sub_type" binding:"required"`
	ResNameEn       string             `json:"res_name_en" form:"res_name_en" binding:"-"`
	ResNameCn       string             `json:"res_name_cn" form:"res_name_en" binding:"required"`
	ResEndpRoute    string             `json:"res_endp_route" form:"res_endp_route" binding:"-"`
	ResDataLocation model.LocationData `json:"res_data_location" form:"res_data_location" binding:"-"`
}
