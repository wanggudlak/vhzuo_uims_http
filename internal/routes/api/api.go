package api

import (
	"github.com/gin-gonic/gin"
	"uims/internal/controllers/auth_controller"
	"uims/internal/controllers/client_controller"
	"uims/internal/controllers/file_controller"
	"uims/internal/controllers/login_controller"
	"uims/internal/controllers/org_controller"
	"uims/internal/controllers/resource_controller"
	"uims/internal/controllers/resource_group_controller"
	"uims/internal/controllers/role_controller"
	"uims/internal/controllers/sms_controller"
	"uims/internal/controllers/user_controller"
	"uims/internal/controllers/user_info_controller"
	"uims/internal/controllers/user_role_controller"
	"uims/internal/controllers/verifycode_controller"
	"uims/internal/controllers/wechatcontroller"
	"uims/internal/middleware"
)

func LoadApi(router *gin.Engine) {
	noAuthAPI := router.Group("/api")
	{
		//public := noAuthAPI.Group("/token")
		//{
		//	public.GET("/findpasswdform", auth_controller.MakeTokenForFindPasswdForm)
		//}
		captcha := noAuthAPI.Group("/captcha")
		{
			captcha.GET("/math", verifycode_controller.GenerateMathCaptchaBase64)
			captcha.GET("/slide", verifycode_controller.GenerateSlideRangeCoordinatePoints)
			captcha.POST("/verifyslide", verifycode_controller.VerifySlideCaptchaLocationAndSendSMS)
		}

		sms := noAuthAPI.Group("/sms/verifycode")
		{
			sms.GET("/send", sms_controller.Send)
			sms.POST("/verify", sms_controller.Verify)
		}

		loginAuthentication := noAuthAPI.Group("/login", middleware.Middleware.CSRF...)
		{
			// 获取登录信息加密公钥
			loginAuthentication.GET("/key", login_controller.GetRSAPubKey)
			// 网页登录
			loginAuthentication.POST("/authenticate", login_controller.Authenticate)
		}

		login := noAuthAPI.Group("/login")
		{
			// 微信网页扫码登录
			login.GET("/wechat/code", login_controller.WeChatCodeLogin)
			// 获取关注二维码
			login.GET("/wechat/qr_code", login_controller.WeChatQRCode)
			// 获取扫码结果
			login.GET("/wechat/qr_code/login", login_controller.WeChatQRCodeLogin)
			// 需要绑定手机号, 则进行绑定
			login.POST("/wechat/qr_code/bind/phone", login_controller.WeChatQRCodeLoginBindPhone)
			// 小程序授权登录
			login.GET("/applet/code", login_controller.AppletCodeLogin)
			login.POST("/applet/info", login_controller.AppletInfo)
		}

		oauth := noAuthAPI.Group("/oauth")
		{
			oauth.GET("wechat/bind", wechatcontroller.Bind)
			oauth.GET("wechat/unbind", wechatcontroller.Unbind)
		}

		auth := noAuthAPI.Group("/auth", middleware.Middleware.CSRF...)
		{
			auth.GET("find/password/token", auth_controller.FindPasswordToken)
			auth.POST("find/password", auth_controller.FindPassword)
			auth.POST("register", auth_controller.Register)
			auth.POST("registerlogin", auth_controller.RegisterAndLogin)
		}

		background_login := noAuthAPI.Group("/background_login")
		{
			background_login.POST("", login_controller.BackgroundLogin)
		}
	}

	api := router.Group("/api", middleware.Middleware.Api...)
	api.Use(middleware.BackJWTMiddleware())
	{
		file := api.Group("/upload")
		{
			file_controller.SetMaxMultipleMemory(router)
			{
				// 上传客户端系统的html模板页
				file.POST("/htmltmpl", file_controller.UploadClientHTMLtmpl)
				file.POST("client/file", file_controller.UploadFile) // 客户端上传KEY文件
			}
		}
		user := api.Group("/users")
		{
			user.POST("", user_controller.Create)
			user.PUT("", user_controller.Update)
			user.GET("", user_controller.Get)

			user.GET("/list", user_controller.List)
			userRole := user.Group("/role")
			{
				userRole.POST("", user_role_controller.AddRole)      // 为用户关联角色
				userRole.DELETE("", user_role_controller.DeleteRole) // 为用户删除角色
			}
		}
		user_info := api.Group("/user/info")
		{
			user_info.POST("", user_info_controller.Create)
			user_info.GET("", user_info_controller.Get)
			user_info.PUT("", user_info_controller.Update)

		}

		resource := api.Group("/resource")
		{
			resource.POST("", resource_controller.Create)
			resource.GET("/list", resource_controller.List)
			resource.DELETE("", resource_controller.Delete)
			resource.PUT("", resource_controller.Update)
			resource.GET("", resource_controller.Get)
		}
		resource_group := api.Group("/resource/group")
		{
			resource_group.POST("", resource_group_controller.Create)
			resource_group.GET("/list", resource_group_controller.List)
			resource_group.DELETE("", resource_group_controller.Delete)
			resource_group.PUT("", resource_group_controller.Update)
			resource_group.GET("", resource_group_controller.Get)
		}

		//role := api.Group("/roles", middleware.Middleware.BackJWT...)
		role := api.Group("/roles")
		{
			role.GET("list", role_controller.List)             //获取角色列表
			role.GET("", role_controller.Detail)               // 获取角色详情信息
			role.POST("", role_controller.AddRole)             // 角色创建
			role.PUT("", role_controller.UpdateRole)           // 角色更新
			role.DELETE("", role_controller.DeleteRole)        // 角色软删除
			role.GET("/group", role_controller.GroupDetail)    // 角色所关联的资源组获取
			role.POST("/group", role_controller.AddGroup)      // 角色添加资源组
			role.DELETE("/group", role_controller.DeleteGroup) // 角色删除资源组
		}
		client := api.Group("/client")
		{
			client.GET("/list", client_controller.List)                  // 获取所有入驻客户端系统
			client.GET("", client_controller.Detail)                     // 获取客户端信息
			client.POST("", client_controller.AddClient)                 //客户端系统基本信息录入
			client.PUT("", client_controller.UpdataClient)               //客户端系统基本信息修改
			client.POST("/status", client_controller.ChangeClientStatus) // 客户端系统 基本信息状态修改
			clientSetting := client.Group("/setting")
			{
				clientSetting.POST("", client_controller.AddClientSetting)      //客户端系统设置信息录入
				clientSetting.PUT("", client_controller.UpdateClientSetting)    //客户端系统设置信息更新
				clientSetting.DELETE("", client_controller.DeleteClientSetting) //客户端系统设置信息删除
			}
		}
		organization := api.Group("/org")
		{
			organization.GET("", org_controller.Detail) // 获取组织详细信息
		}

	}
}
