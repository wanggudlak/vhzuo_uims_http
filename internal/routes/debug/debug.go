package debug

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"uims/conf"
	"uims/internal/service"
	"uims/pkg/db"
	"uims/pkg/gredis"
	thriftclient "uims/pkg/thrift/client"
)

func LoadDebug(r *gin.Engine) {
	debug := r.Group("debug")
	debug.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	debug.GET("/health", func(c *gin.Context) {
		// 数据库状态
		var mysql = gin.H{}
		for name, config := range conf.Database.MySQL {
			mysql[name] = gin.H{
				"host":     config.Host + ":" + config.Port,
				"database": config.Database,
				"user":     config.Username,
				"ok":       db.Conn(name).DB().Ping() == nil,
				"err":      db.Conn(name).DB().Ping(),
			}
		}
		// redis 状态
		var redisStatus = gin.H{}
		for name, config := range conf.Database.Redis {
			redisStatus[name] = gin.H{
				"host":     fmt.Sprintf("%s:%d", config.Host, config.Port),
				"database": config.Database,
				"ok":       gredis.Conn(name).Ping().Val() == "PONG",
				"err":      gredis.Conn(name).Ping().Err(),
			}
		}
		// thrift server 状态
		// thrift client 状态
		var thriftClients = gin.H{}
		for name, config := range conf.ThriftClients {
			resp := service.Response{}
			err := func() error {
				cli, err := thriftclient.Get(&config)
				if err != nil {
					return err
				}
				fmt.Printf("call %s \n", config.ServerAddr)
				err = cli.Call(service.Request{
					BRequest: thriftclient.BRequest{
						MethodName: "debug_call",
						Params:     nil,
					},
				}, &resp)
				if err != nil {
					return err
				}
				fmt.Printf("%+v \n", resp)
				return nil
			}()
			if err != nil {
				thriftClients[name] = gin.H{
					"config": config,
					"ok":     false,
					"err":    err.Error(),
				}
			} else {
				thriftClients[name] = gin.H{
					"config": config,
					"ok":     true,
					"err":    nil,
					"resp":   resp,
				}
			}
		}
		// 域名配置状态
		// 开关状态
		c.JSON(http.StatusOK, gin.H{
			"name":           conf.Name,
			"env":            conf.Env,
			"debug":          conf.Debug,
			"host":           conf.Host,
			"redis":          redisStatus,
			"mysql":          mysql,
			"switchs":        conf.Switch,
			"thrift_clients": thriftClients,
		})
	})
}
