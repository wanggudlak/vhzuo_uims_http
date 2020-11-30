package conf

import "uims/pkg/env"

const (
	SMS_ALI_DRIVER                  = "aliyun"
	SMS_ALI_SIGN                    = "微桌"
	SMS_ALI_VERIFY_CODE_TEMPLATE_ID = "SMS_169101455"
)

type sms struct {
	smsdriver string
	smsks     interface{}
}

type smsAli struct {
	smsAliRegionID     string
	smsAliAccessKey    string
	smsAliAccessSecret string
}

func newDriverParam(driverType string) interface{} {
	switch driverType {
	case SMS_ALI_DRIVER:
		return &smsAli{
			smsAliRegionID:     env.DefaultGet("SMS_ALI_REGION_ID", "cn-hangzhou").(string),
			smsAliAccessKey:    env.DefaultGet("SMS_ALI_ACCESS_KEY", "").(string),
			smsAliAccessSecret: env.DefaultGet("SMS_ALI_ACCESS_SECRET", "").(string),
		}
	default:
		panic("SMS driver type is invalid.")
	}
}

func NewSMSconf() *sms {
	driverType := env.DefaultGet("SMS_DRIVER", "aliyun").(string)
	return &sms{
		smsdriver: driverType,
		smsks:     newDriverParam(driverType),
	}
}

func (smsConf *sms) GetDriverType() string {
	return smsConf.smsdriver
}

func (smsConf *sms) GetDriverParam() interface{} {
	return smsConf.smsks
}

func (smsConf *sms) GetAliDriverParam() *smsAli {
	return smsConf.GetDriverParam().(*smsAli)
}

func (smsAliConf *smsAli) GetSMSaliRegionID() string {
	return smsAliConf.smsAliRegionID
}

func (smsAliConf *smsAli) GetSMSaliAccessKey() string {
	return smsAliConf.smsAliAccessKey
}

func (smsAliConf *smsAli) GetSMSaliAccessSecret() string {
	return smsAliConf.smsAliAccessSecret
}
