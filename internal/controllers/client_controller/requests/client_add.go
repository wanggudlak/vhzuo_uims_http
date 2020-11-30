package requests

type ClientNewRequest struct {
	Type     string `json:"type" form:"type" binding:"required" comment:"客户端类型，VDK：微桌"`
	FlagCode string `json:"flag_code" form:"flag_code" binding:"required" comment:"客户端业务系统标识，VDK_CASS：微桌结算系统等"`
	Spm1Code string `json:"spm1_code" form:"spm1_code" binding:"required" comment:"SPM编码中的第一部分，微桌内部系统用1024；外部系统用2048"`
	Spm2Code string `json:"spm2_code" form:"spm2_code" binding:"-" comment:"SPM编码中的第二部分"`
	Name     string `json:"name" form:"name" binding:"required" comment:"客户端业务系统名称"`
	HostIP   string `json:"host_ip" form:"host_ip" binding:"required" comment:"客户端当前使用的IP，多个用json字符串保存"`
	HostURL  string `json:"host_url" form:"host_url" binding:"required" comment:"客户端业务系统当前使用的域名，多个用json字符串保存"`
	INAT     string `json:"in_at" form:"in_at" binding:"-" comment:"入驻可以使用的开始时间点，默认为当前时间"`
	ForgetAT string `json:"forget_at" form:"forget_at" binding:"-"  comment:"客户端系统使用UIMS失效时间"`
}
