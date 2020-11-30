package uuid

import (
	"fmt"
	"uims/internal/service/uuid"
	thriftserver "uims/pkg/thrift/server"
)

func Generate(c *thriftserver.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Response.Error(fmt.Errorf("Generate UUID failed: %s", err.(error).Error()))
		}
	}()

	id := uuid.GenerateForCASS()

	c.Response.Success(id, "success")
}
