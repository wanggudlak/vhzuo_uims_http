package conf

import (
	gin2 "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
	"uims/app"
	_ "uims/app"
	"uims/pkg/env"
	"uims/pkg/glog/hook"
	thriftclient "uims/pkg/thrift/client"
)

var (
	GinModel  string
	Name      = os.Getenv("APP_NAME")
	URL       = os.Getenv("APP_URL")
	Env       = os.Getenv("APP_ENV")
	APPKey    = os.Getenv("APP_KEY")
	Debug     = env.DefaultGetBool("DEBUG", false)
	Host      = os.Getenv("APP_HOST")
	DebugHost = env.DefGetStr("APP_DEBUG_HOST", "127.0.0.1:12548")
	Database  = database{
		MySQL: map[string]MysqlConf{
			"default": {
				Host:        os.Getenv("DB_HOST"),
				Port:        os.Getenv("DB_PORT"),
				Username:    os.Getenv("DB_USERNAME"),
				Password:    os.Getenv("DB_PASSWORD"),
				Database:    os.Getenv("DB_DATABASE"),
				MaxLiftTime: time.Second * 60,
			},
			"cass": {
				Host:        os.Getenv("CASS_DB_HOST"),
				Port:        os.Getenv("CASS_DB_PORT"),
				Username:    os.Getenv("CASS_DB_USERNAME"),
				Password:    os.Getenv("CASS_DB_PASSWORD"),
				Database:    os.Getenv("CASS_DB_DATABASE"),
				MaxLiftTime: time.Second * 60,
			},
			//"task": {
			//	Host:        os.Getenv("TASK_DB_HOST"),
			//	Port:        os.Getenv("TASK_DB_PORT"),
			//	Username:    os.Getenv("TASK_DB_USERNAME"),
			//	Password:    os.Getenv("TASK_DB_PASSWORD"),
			//	Database:    os.Getenv("TASK_DB_DATABASE"),
			//	MaxLiftTime: time.Second * 60,
			//},
		},
		Redis: map[string]RedisConf{
			"default": {
				Host:     env.DefaultGet("REDIS_HOST", "127.0.0.1").(string),
				Password: env.DefaultGet("REDIS_PASSWORD", "").(string),
				Port:     env.DefaultGetInt("REDIS_PORT", 6379),
				Database: env.DefaultGetInt("REDIS_DATABASE", 0),
			},
		},
	}
	Filesystems = filesystems{
		Default: "local",
		Cloud:   "",
		Disks: Disks{
			Local: Local{
				Driver: "local",
				Root:   "app/public",
			},
		},
	}
	Logging = struct {
		Channels logs
		Default  string
	}{
		Channels: logs{
			"default": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/def/main.log"),
				Level:  log.DebugLevel,
				Days:   7,
				Hooks: []log.Hook{
					&hook.DefaultFieldHook{
						AppName: Name,
						AppUrl:  URL,
						AppEnv:  Env,
					},
				},
			},
			"gin": Log{
				Driver:       env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:         app.GetStoragePath("logs/route/route.log"),
				Level:        log.DebugLevel,
				Days:         7,
				LogFormatter: &log.TextFormatter{},
			},
			"db": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/db/db.log"),
				Level:  log.DebugLevel,
				Days:   7,
			},
			"cron": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/cron/cron.log"),
				Level:  log.DebugLevel,
				Days:   7,
			},
			"thrift": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/thrift/thrift.log"),
				Level:  log.DebugLevel,
				Days:   7,
			},
			"thriftclient": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/thriftclient/thriftclient.log"),
				Level:  log.DebugLevel,
				Days:   7,
			},
			"sms": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/sms/sms.log"),
				Level:  log.DebugLevel,
				Days:   30,
			},
			"request": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/request/request.log"),
				Level:  log.DebugLevel,
				Days:   30,
			},
			"casswechat": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/casswechat/casswechat.log"),
				Level:  log.DebugLevel,
				Days:   30,
			},
			"applet": Log{
				Driver: env.DefGetStr("LOG_DEFAULT_DRIVER", Daily),
				Path:   app.GetStoragePath("logs/applet/applet.log"),
				Level:  log.DebugLevel,
				Days:   30,
			},
		},
		Default: "default",
	}
	SMS        = NewSMSconf()
	AppSetting = setting{
		PageSize:  env.DefaultGetInt("PAGE_SIZE", 10),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
	Switch = SwitchControl{
		ImgCaptcha: env.DefaultGetBool("SWITCH_IMG_CAPTCHA", true),
		SMSCaptcha: env.DefaultGetBool("SWITCH_SMS_CAPTCHA", true),
		CSRF:       env.DefaultGetBool("SWITCH_CSRF", true),
	}
	EmailConf     = NewEmailConf()
	UUIDConf      = NewNodeMap()
	ThriftClients = map[string]thriftclient.Config{
		"cass": {
			OnOff:                  env.DefaultGetBool("THRIFT_CLIENT_ON_OFF_CASS", true),
			ServerAddr:             env.DefaultGet("THRIFT_CLIENT_SERVER_ADDR_CASS", "127.0.0.1:9091").(string),
			DataProtocol:           env.DefGetStr("THRIFT_CLIENT_PROTOCOL_CASS", "binary"),
			BufferedSize:           env.DefaultGetInt("THRIFT_CLIENT_BUFFERED_SIZE_CASS", 8192),
			Buffered:               env.DefaultGetBool("THRIFT_CLIENT_BUFFERED_CASS", false),
			Framed:                 env.DefaultGetBool("THRIFT_CLIENT_FRAMED_CASS", true),
			Secure:                 env.DefaultGetBool("THRIFT_CLIENT_SECURE_CASS", true),
			IsUseIOMultiplexing:    env.DefaultGetBool("THRIFT_CLIENT_USE_IOMULTIPLEXING_CASS", true),
			ServerAPIServiceLoc:    env.DefaultGet("THRIFT_CLIENT_SERVER_API_SERVICE_LOC", "UIMSRpcApiService").(string),
			InitialConnCountInPool: 5,
			MaxConnCountOfPool:     30,
			SocketTimeout:          time.Second * 60,
		},
		"mp": {
			OnOff:                  env.DefaultGetBool("THRIFT_CLIENT_ON_OFF_MP", true),
			ServerAddr:             env.DefaultGet("THRIFT_CLIENT_SERVER_ADDR_MP", "127.0.0.1:9091").(string),
			DataProtocol:           env.DefGetStr("THRIFT_CLIENT_PROTOCOL_MP", "binary"),
			BufferedSize:           env.DefaultGetInt("THRIFT_CLIENT_BUFFERED_SIZE_MP", 8192),
			Buffered:               env.DefaultGetBool("THRIFT_CLIENT_BUFFERED_MP", false),
			Framed:                 env.DefaultGetBool("THRIFT_CLIENT_FRAMED_MP", true),
			Secure:                 env.DefaultGetBool("THRIFT_CLIENT_SECURE_MP", true),
			IsUseIOMultiplexing:    env.DefaultGetBool("THRIFT_CLIENT_USE_IOMULTIPLEXING_MP", true),
			ServerAPIServiceLoc:    env.DefaultGet("THRIFT_CLIENT_SERVER_API_SERVICE_LOC", "UIMSRpcApiService").(string),
			InitialConnCountInPool: 5,
			MaxConnCountOfPool:     30,
			SocketTimeout:          time.Second * 60,
		},
	}
)

func init() {
	if !strings.EqualFold(Env, "local") &&
		!strings.EqualFold(Env, "production") &&
		!strings.EqualFold(Env, "testing") {
		panic("env APP_ENV must be: local, production, testing")
	}
	switch Env {
	case "testing":
		GinModel = gin2.TestMode
	case "local":
		GinModel = gin2.DebugMode
	case "production":
		GinModel = gin2.ReleaseMode
	}
}
