package boot

import (
	"log"
	"uims/app"
	"uims/internal/middleware"
	"uims/internal/routes"
	"uims/pkg/db"
	"uims/pkg/glog"
	"uims/pkg/gredis"
	migrate2 "uims/pkg/migrate"
	"uims/pkg/storage"
)

func SetInTest() {
	app.InTest = true
}

func SetInConsole() {
	app.InConsole = true
}

// 应用启动入口
func Boot() {
	var err error
	glog.Init()
	storage.Init(app.StoragePath)

	if _, err = gredis.InitDef(); err != nil {
		log.Panicf("Init Default Redis connection filed: %+v", err)
	}

	if _, err = db.InitDef(); err != nil {
		log.Panicf("Init Default MySQL connection filed: %+v", err)
	}

	// 命令行模式下不加载路由
	if !app.InConsole {
		// 注册中间件
		middleware.Init()
		// 注册路由
		router := routes.InitRouter()
		app.SetEngineRouter(router)
	}

	//validator.Init()

	app.Booted = true

	migrate()
}

func Destroy() {
	db.Close()
	glog.Close()
	gredis.Close()
}

func migrate() {
	db.Def().AutoMigrate(&migrate2.Migration{})
}
