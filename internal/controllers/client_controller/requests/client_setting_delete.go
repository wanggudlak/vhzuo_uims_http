package requests

type ClientSettingDeleteRequest struct {
	ID int `json:"id" form:"id" binding:"required" comment:"客户设置信息ID"`
}
