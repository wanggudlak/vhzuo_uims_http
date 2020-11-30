package democontroller

import thriftserver "uims/pkg/thrift/server"

type DemoReq struct {
	A int `json:"a" binding:"required"`
	B int `json:"b" binding:"required"`
	C int `json:"c" binding:"required"`
}

func Demo(c *thriftserver.Context) {
	var req DemoReq
	if err := c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	// .. 调用具体业务逻辑方法
	c.Response.Success(req, "业务处理完成")
	return
}
