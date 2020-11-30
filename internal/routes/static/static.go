package static

import (
	"github.com/gin-gonic/gin"
	"uims/internal/controllers/login_controller/html_controller"
	"uims/internal/service"
)

// LoadStaticRouter 注册静态文件访问所需的路由
func LoadStaticRouter(router *gin.Engine) {
	err := service.GetAllClientsNeedRenderHTML(&service.AllClientsNeedRenderHTML, []string{"*"})
	if err != nil {
		panic(err)
	}

	if len(service.AllClientsNeedRenderHTML) > 0 {
		loadedStaticRoute := map[string]bool{}
		loadedStaticAPIRoute := map[string]bool{}
		for _, clientSetting := range service.AllClientsNeedRenderHTML {
			staticRouteURIPrefix := service.GetStaticFileRouteURIPrefix()
			staticRouteURI := service.GetStaticFileRouteURI(clientSetting.SpmFullCode)
			staticRouteSuffix := service.GetStaticFileRouteSuffix(&clientSetting)
			staticFileStoragePath := service.GetStaticFileStoragePath(staticRouteURI)

			if loaded, ok := loadedStaticRoute[staticRouteURI]; ok && loaded {
				continue
			} else {
				router.Static(staticRouteURI, staticFileStoragePath)
				loadedStaticRoute[staticRouteURI] = true

				if len(staticRouteSuffix) > 0 {
					if loaded2, ok2 := loadedStaticAPIRoute[staticRouteSuffix]; ok2 && loaded2 {
						continue
					} else {
						noAuthGetHTML := router.Group(staticRouteURIPrefix)
						{
							noAuthGetHTML.GET(staticRouteSuffix, html_controller.RenderHTML)
						}
						loadedStaticAPIRoute[staticRouteSuffix] = true
					}
				}
			}
		}
	}
}
