package requests

type UserInfoUpdateRequest struct {
	Birthday              string `form:"birthday" json:"birthday" binding:"required"`
	HeaderImgURL          string `form:"header_img_url" json:"header_img_url" binding:"required"`
	IdentityCardEmblemImg string `form:"identity_card_emblem_img" json:"identity_card_emblem_img" binding:"-"`
	IdentityCardNo        string `form:"identity_card_no" json:"identity_card_no" binding:"-"`
	IdentityCardPersonImg string `form:"identity_card_person_img" json:"identity_card_person_img" binding:"-"`
	LandlinePhone         string `form:"landline_phone" json:"landline_phone" binding:"-"`
	Nickname              string `form:"nickname" json:"nickname" binding:"required"`
	Sex                   string `form:"sex" json:"sex" binding:"required"`
	TaxerNo               string `form:"taxer_no" json:"taxer_no" binding:"-"`
	TaxerType             string `form:"taxer_type" json:"taxer_type" binding:"-"`
	UserID                int    `form:"user_id" json:"user_id" binding:"required"`
	Wechat                string `form:"wechat" json:"wechat" binding:"-"`
	NameCn                string `form:"name_cn" json:"name_cn" binding:"-"`
	NameEn                string `form:"name_en" json:"name_en" binding:"-"`
	NameCnAlias           string `form:"name_cn_alias" json:"name_cn_alias" binding:"-"`
}
