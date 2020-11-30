package responses

type SMSaliResponse struct {
	Message   string
	RequestId string
	Code      string
}

func (smsAliResp *SMSaliResponse) IsSendSuccess() (result bool, message string) {
	result = smsAliResp.Code == "OK"
	message = smsAliResp.Message
	return
}
