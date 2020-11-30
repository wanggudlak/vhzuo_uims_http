package requests

type ClientUpdateRequest struct {
	ID             int    `json:"id" form:"id" binding:"required" comment:"客户端类型，VDK：微桌"`
	Type           string `json:"type" form:"type" binding:"-" comment:"客户端类型，VDK：微桌"`
	FlagCode       string `json:"flag_code" form:"flag_code" binding:"-" comment:"客户端业务系统标识，VDK_CASS：微桌结算系统等"`
	HostIP         string `json:"host_ip" form:"host_ip" binding:"-" comment:"客户端当前使用的IP，多个用json字符串保存"`
	HostURL        string `json:"host_url" form:"host_url" binding:"-" comment:"客户端业务系统当前使用的域名，多个用json字符串保存"`
	PUBKryPath     string `json:"client_pub_key_path" form:"client_pub_key_path" binding:"-" comment:"客户端业务系统的RSA公钥key文件路径"`
	UIMSPubKeyPath string `json:"uims_pub_key_path" form:"uims_pub_key_path" binding:"-" comment:"UIMS系统的RSA公钥文件路径"`
	UIMSPriKeyPath string `json:"uims_pri_key_path" form:"uims_pri_key_path" binding:"-" comment:"UIMS系统的RSA私钥文件路径"`
	INAT           string `json:"in_at" form:"in_at" binding:"-" comment:"入驻可以使用的开始时间点，默认为当前时间"`
	ForgetAT       string `json:"forget_at" form:"forget_at" binding:"-"  comment:"客户端系统使用UIMS失效时间"`
}
