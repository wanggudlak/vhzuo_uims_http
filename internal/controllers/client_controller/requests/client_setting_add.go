package requests

type NewClientSettingRequest struct {
	ClientID int    `json:"client_id" form:"client_id" binding:"required" comment:"客户端业务ID"`
	Type     string `json:"type" form:"type" binding:"required" comment:"类型：LGN-用于登录的设置；REG-用于注册的设置；"`
	//Fields       map[string]interface{} `json:"form_fields" form:"form_fields" binding:"-" comment:"表单域属性数据，内容用json进行存储"`
	//TemplateFile model.PageTemplateFile `json:"page_template_file" form:"page_template_file" binding:"required" comment:"登录页或注册页html模板文件路径"`
	TemplateFile string `json:"page_template_file" form:"page_template_file" binding:"required" comment:"登录页或注册页html模板文件路径"`
}
