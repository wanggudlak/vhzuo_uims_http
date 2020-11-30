package model

import (
	"time"
)

type Client struct {
	ID               uint      `json:"id" gorm:"primary_key;comment:'用户或会员ID'"`
	AppId            string    `json:"app_id" gorm:"column:app_id;type:char(16);unique;not null;default:'';comment:'客户端系统APPID，用来唯一标识客户端系统'"`
	AppSecret        string    `json:"app_secret" gorm:"column:app_secret;type:varchar(255);not null;default:'';comment:'客户端系统与UIMS系统之间用来对称加解密的秘钥，base64编码后存储'"`
	ClientType       string    `json:"type" gorm:"column:client_type;type:char(32);default:'';comment:'客户端类型，VDK：微桌，CASS：结算系统'"`
	ClientFlagCode   string    `json:"flag_code" gorm:"column:client_flag_code;type:varchar(16);default:'';comment:'客户端业务系统标识，VDK_MP：微桌任务系统平台；VDK_CRM：微桌CRM系统；VDK_INVO：微桌代开发票系统；VDK_ESIGN：微桌电签系统；VDK_ES_SAPP：微桌电签小程序；VDK_CASS_FRONT：结算系统前台；VDK_CASS_BACK：结算系统后台'"`
	ClientSpm1Code   string    `json:"spm1_code" gorm:"column:client_spm1_code;type:char(4);default:'1024';comment:'客户端业务系统的SPM编码中的第一部分（外站类型ID），微桌内部系统用1024；外部系统用2048'"`
	ClientSpm2Code   string    `json:"spm2_code" gorm:"column:client_spm2_code;type:char(16);default:'';comment:'客户端业务系统的SPM编码中的第二部分（外站APP ID），和APPID一致'"`
	ClientName       string    `json:"name" gorm:"column:client_name;type:varchar(255);default:'';comment:'客户端业务系统名称'"`
	Status           string    `json:"status" gorm:"column:status;type:char(1);default:'N';comment:'客户端业务系统使用UIMS的状态，默认N：未授权不可用；Y：已授权可用；F-被禁用'"`
	ClientHostIp     string    `json:"host_ip" gorm:"column:client_host_ip;type:varchar(255);default:''comment:'客户端使用的IP，可用于白名单'"`
	ClientHostUrl    string    `json:"host_url" gorm:"column:client_host_url;type:varchar(255);default:''comment:'客户端业务系统当前使用的域名，例如微桌结算系统是https://fuwu.skysharing.cn'"`
	ClientPubKeyPath string    `json:"client_pub_key_path" gorm:"column:client_pub_key_path;type:varchar(128);default:'';comment:'客户端业务系统的RSA公钥key文件路径，以appid作为各自的目录'"`
	UIMSPubKeyPath   string    `json:"uims_pub_key_path" gorm:"column:uims_pub_key_path;type:varchar(128);default:'';comment:'UIMS系统的RSA公钥文件路径'"`
	UIMSPriKeyPath   string    `json:"uims_pri_key_path" gorm:"column:uims_pri_key_path;type:varchar(128);default:'';comment:'UIMS系统的RSA私钥文件路径'"`
	InAt             time.Time `json:"in_at" gorm:"column:in_at;type:datetime;default:null;comment:'入驻可以使用的开始时间点'"`
	ForgetAt         time.Time `json:"forget_at" gorm:"column:forget_at;type:datetime;default:null;comment:'在什么时间点，客户端系统不能使用UIMS，默认是空字符串'"`
	*CommonModel
}

func (Client) TableName() string {
	return "uims_client"
}
