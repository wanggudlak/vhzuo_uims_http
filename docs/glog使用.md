# glog 使用

## 提供功能

- 支持 daily, single 驱动
- 可增加渠道
- 可配置参数
- 可分渠道
- 可按日期划分
- 可按规则自动清理

## 配置

`vim conf/config.go`
```vim
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
		},
		Default: "default",
	}
```

## 使用

```
package main

import (
 "uims/pkg/glog"
 log "github.com/sirupsen/logrus"
)

// 全局使用
log.Print("test")
// 与全局使用一样使用 default 驱动
glog.Default().Print("test")
glog.Channel("gin").Print("gin test")

```