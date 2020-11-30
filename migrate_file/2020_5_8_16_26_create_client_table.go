package migrate_file

import (
	"fmt"
	"time"
	"uims/internal/model"
	"uims/pkg/db"
)

type Client struct {
	ID               uint      `gorm:"primary_key;comment:'用户或会员ID'"`
	AppId            string    `gorm:"column:app_id;type:char(16);unique;not null;default:'';comment:'客户端系统APPID，用来唯一标识客户端系统'"`
	AppSecret        string    `gorm:"column:app_secret;type:varchar(255);not null;default:'';comment:'客户端系统与UIMS系统之间用来对称加解密的秘钥，base64编码后存储'"`
	ClientType       string    `gorm:"column:client_type;type:char(32);not null;default:'';comment:'客户端类型，VDK：微桌，CASS：结算系统'"`
	ClientFlagCode   string    `gorm:"column:client_flag_code;type:varchar(16);not null;default:'';comment:'客户端业务系统标识，VDK_CASS：微桌结算系统；VDK_MP：微桌任务系统平台；VDK_CRM：微桌CRM系统；VDK_INVO：微桌代开发票系统；VDK_ESIGN：微桌电签系统；VDK_ES_SAPP：微桌电签小程序；'"`
	ClientSpm1Code   string    `gorm:"column:client_spm1_code;type:char(4);not null;default:'1024';comment:'客户端业务系统的SPM编码中的第一部分（外站类型ID），微桌内部系统用1024；外部系统用2048'"`
	ClientSpm2Code   string    `gorm:"column:client_spm2_code;type:char(16);not null;default:'';comment:'客户端业务系统的SPM编码中的第二部分（外站APP ID），和APPID一致'"`
	ClientName       string    `gorm:"column:client_name;type:varchar(255);not null;default:'';comment:'客户端业务系统名称'"`
	Status           string    `gorm:"column:status;type:char(1);not null;default:'N';comment:'客户端业务系统使用UIMS的状态，默认N：未授权不可用；Y：已授权可用；F-被禁用'"`
	ClientHostIp     string    `gorm:"column:client_host_ip;type:varchar(255);comment:'客户端使用的IP，可用于白名单'"`
	ClientHostUrl    string    `gorm:"column:client_host_url;type:varchar(255);comment:'客户端业务系统当前使用的域名，例如微桌结算系统是https://fuwu.skysharing.cn'"`
	ClientPubKeyPath string    `gorm:"column:client_pub_key_path;type:varchar(128);not null;default:'';comment:'客户端业务系统的RSA公钥key文件路径，以appid作为各自的目录'"`
	UIMSPubKeyPath   string    `gorm:"column:uims_pub_key_path;type:varchar(128);not null;default:'';comment:'UIMS系统的RSA公钥文件路径'"`
	UIMSPriKeyPath   string    `gorm:"column:uims_pri_key_path;type:varchar(128);not null;default:'';comment:'UIMS系统的RSA私钥文件路径'"`
	InAt             time.Time `gorm:"column:in_at;comment:'入驻可以使用的开始时间点'"`
	ForgetAt         time.Time `gorm:"column:forget_at;comment:'在什么时间点，客户端系统不能使用UIMS，默认是空字符串'"`
	*model.CommonModel
}

func (Client) TableName() string {
	return "uims_client"
}

type CreateClientTableMigrate struct {
}

func (CreateClientTableMigrate) Key() string {
	return "2020_5_8_16_26_create_client_table"
}

func (CreateClientTableMigrate) Up() (err error) {
	if db.Def().HasTable(Client{}.TableName()) {
		err = fmt.Errorf("users table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='客户端即使用uims系统的业务系统'").
		CreateTable(&Client{}).Error
	return
}

func (CreateClientTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&Client{}).Error
	return
}
