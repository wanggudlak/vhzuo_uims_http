package routes

import (
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strings"
	"uims/app"
	"uims/conf"
	"uims/internal/middleware"
	"uims/internal/routes/api"
	"uims/internal/routes/debug"
	"uims/internal/routes/static"
	"uims/internal/routes/swag"
	"uims/pkg/glog"
)

func InitRouter() *gin.Engine {
	var router *gin.Engine
	gin.SetMode(conf.GinModel)
	if app.InTest || app.InConsole {
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(glog.Channel("gin").Out)
	} else {
		// 非测试或命令将输出路由信息到屏幕上
		gin.ForceConsoleColor()
		gin.DefaultWriter = io.MultiWriter(glog.Channel("gin").Out, os.Stdout)
	}

	router = gin.New()
	router.Use(gin.RecoveryWithWriter(glog.Channel("gin").Out), gin.Logger())
	// 加载默认中间件
	router.Use(middleware.Middleware.Def...)
	loadRoutes(router)
	return router
}

// 新增加的路由文件需要在这里进行加载
func loadRoutes(router *gin.Engine) {
	router.Use(middleware.Cors())
	// 注册静态文件访问所需的路由
	static.LoadStaticRouter(router)

	// 注册请求API所需的路由
	api.LoadApi(router)

	if strings.EqualFold(conf.Env, "local") || strings.EqualFold(conf.Env, "testing") {
		// 注册访问swagger文档所需的路由
		swag.LoadSwag(router)
	}
}

func InitDebugRouter() *gin.Engine {
	var router *gin.Engine
	router = gin.Default()
	debug.LoadDebug(router)
	return router
}
