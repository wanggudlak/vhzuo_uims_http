package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"uims/conf"
	resp "uims/internal/controllers/responses"
	"uims/internal/service"
	"uims/internal/service/adapter"
)

type Request struct {
	Token    string `json:"_token" form:"_token" binding:"required" comment:"CSRF-TOKEN"`
	ClientID int    `json:"client_id" form:"client_id" binding:"required" comment:"业务系统ID"`
}

func VerifyCSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.Switch.CSRF {
			c.Next()
			return
		}
		var err error
		var req Request
		err = func() error {
			// 从原有Request.Body读取
			//bodyBytes, err = ioutil.ReadAll(c.Request.Body)
			method := c.Request.Method
			if "GET" != method {
				err = adapter.New().ReadRequestBody(c.Request, &req)
				if err != nil {
					return fmt.Errorf("invalid request body")
				}
			}

			if len(req.Token) == 0 {
				req.Token = c.DefaultQuery("_token", "")
				if len(req.Token) == 0 {
					return fmt.Errorf("请求未携带csrf-token，请重新刷新页面")
				}
			}
			// 取出客户端id
			if 0 == req.ClientID {
				req.ClientID, err = strconv.Atoi(c.DefaultQuery("client_id", ""))
				if 0 == req.ClientID {
					return fmt.Errorf("由于您的csrf-token已过期失效或者页面已被篡改，请重新刷新页面")
				}
			}

			// 验证csrf-token
			if !service.VerifyCSRFToken(req.Token) {
				return service.CSRFTokenVerifiedFailedErr
			}
			return nil
		}()

		if err != nil {
			c.Abort()
			resp.Error(c, err)
			return
		}
		c.Next()
	}
}
