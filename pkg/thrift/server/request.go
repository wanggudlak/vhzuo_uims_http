package thriftserver

type Request struct {
	// 请求方法名
	MethodName string
	// 原始字符串数据
	Body string
}

func NewRequest() *Request {
	return &Request{}
}
