# 目录说明

1. 业务相关代码目录放在 internal 目录中
2. 业务无关代码放在 pkg 目录中
3. pkg 中的包之间尽量不要互相依赖

.
├── Makefile make 命令控制软件生命周期
├── README.md
├── app
│   ├──app.go 公用数据存放
├── bin 二进制可执行文件编译结果存放目录, 不提交git
├── boot
│   └── boot.go 项目加载起始点
├── command 命令行命令目录
│   ├── cmd_template
│   ├── command.go
│   ├── commands
│   ├── show_users_command.go
│   └── tinker.go
├── conf 配置目录
│   ├── config.go
│   ├── config_test.go
│   ├── database.go
│   ├── filesystems.go
│   ├── logging.go
│   ├── setting.go
│   └── sms.go
├── docs 项目文档目录
│   ├── Apis 接口文档
│   ├── docs.go swagger
│   ├── swagger.json
│   ├── swagger.yaml
│   ├── uims.sql
├── go.mod
├── internal 业务代码
│   ├── controllers 控制器
│   ├── middleware 中间件
│   ├── model 模型
│   ├── routes 路由
│   ├── service 服务
│   └── validator 验证器
├── main.go 项目入口
├── migrate_file 数据库迁移文件
├── pkg 业务无关代码
│   ├── color
│   ├── const_definition
│   ├── db  数据库连接
│   ├── e 错误常量
│   ├── encryption
│   ├── env 环境变量加载
│   ├── migrate 数据库迁移
│   ├── randc 随机值产生
│   ├── redis
│   ├── tool 工具
│   ├── log 日志
│   └── type 类型
├── storage 储存目录
│   ├── app 上传与生成文件存放
│   │   └── public 外部可访问目录
│   ├── logs 日志存放
└── vendor 第三方依赖包
