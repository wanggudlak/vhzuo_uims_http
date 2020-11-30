package common

const (
	COMPACT     = "compact"
	SIMPLE_JSON = "simplejson"
	JSON        = "json"
	BINARY      = "binary"
)

const KEY_DIR = "./key/thrift_server"
const PRIVATE_KEY_PASSWD = "uims"

const LOGGER_CHANNEL_TAG = "thrift"

//const METHOD_NAME = "method_name"
//const PARAMS = "params"
const (
	// 请求成功时响应 msg
	//CALL_FAILED_MSG     = "Invoke failed"

	CALL_SUCCESS_MSG    = "Invoke success"
	INVALID_METHOD_NAME = "请求参数错误：未解析出所要请求的接口名"
	INVALID_PARAMS      = "请求参数错误：解析请求参数失败"
)

// STATUS
const (
	STATUS_SUCCESS       = "success"
	STATUS_FAILED        = "failed"
	STATUS_PARAMS_FAILED = "params_failed"
)
