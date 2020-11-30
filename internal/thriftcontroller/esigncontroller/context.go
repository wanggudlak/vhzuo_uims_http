package esigncontroller

type NotifyESignReq struct {
	Name                        string `json:"name" binding:"required"`
	IdCard                      string `json:"id_card" binding:"required"`
	Phone                       string `json:"phone" binding:"omitempty"`
	IdentityCardPersonImgBase64 string `json:"identity_card_person_img_base64" binding:"omitempty"`
	IdentityCardEmblemImgBase64 string `json:"identity_card_emblem_img_base64" binding:"omitempty"`
}

type NotifyESignResp struct {
}
