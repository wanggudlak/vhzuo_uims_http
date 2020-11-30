# uims

## 介绍
统一权限认证管理系统的后端仓库及相关文档

## 目录说明

[目录说明](docs/tree.md)
[glog 日志使用](docs/glog使用.md)

## 软件架构
软件架构说明

### 如何运行
### 本系统用make管理项目的部署
### 安装 GNU make 
```
Mac :  brew install make
Ubuntu: sudo apt-get install make
Other linux : 
    wget http://alpha.gnu.org/gnu/make/  
    cd makefile 文件所在目录
    make 
    make install
```

### 注意：本系统的根目录下有一个makefile文件，用来定义项目部署的make规则，请不要将缩进用的tab格式化为空格！！！

#### 加载依赖

1. 通过 module 模式加载
    - `Goland` 配置 `Go->Go Modules`
        - 勾选 `Enable Go Moudles`
        - `Proxy: direct`
        - `Vendoring mode`
2. 通过 GOPATH 加载 (`待补充...`)

> 加载私有库报错时尝试以下设置

1. 编辑 git 配置文件

    `vim ~/.gitconfig`

    ```config
    [url "git@gitee.com:"]
        insteadOf = https://gitee.com/
    ```

2. 关闭代理加载依赖

    ```shell script
    # 关闭代理
    export GOPROXY=
    # 加载依赖
    go mod tidy
    go mod download
    ```



#### 配置 `.env`

```shell script
cp .env.example .env
# 主要配置以下项
# APP_HOST 服务监听的地址端口
# DB_* 连接到的数据库
# REDIS_* 连接到的 Redis
# DEBUG 是否开启 DEBUG 模式, 主要影响日志的输出
# APP_ENV 配置当前环境, 本地使用 local, 部署使用 production
```

#### 调试

```shell script
如果你的系统上安装了 GNU make,建议通过make管理项目
仅编译:make 或者 make build
编译并开启API server : make run

由于加了私有仓库的包,可以尝试运行以下命令:
构建uims
GOPRIVATE=gitee.com/skysharing go build -o uims main.go
测试例子:
GOPRIVATE=gitee.com/skysharing go test -run=TestGenerateCaptchaBase64
```

#### 同环境编译

```shell script
make 或者 make build
```

#### 显示uims命令行
```
./uims
```

#### 交叉编译

Mac 下编译 Linux 和 Windows 64位可执行程序
```
make linux-in-mac
make win64-in-mac
```
Linux 下编译 Mac 和 Windows 64位可执行程序
```
make mac-in-linux
make win64-in-linux
```

Windows 下编译 Mac 和 Linux 64位可执行程序
```
make mac-in-win64
make linux-in-win64
```

### Thrift RPC 
- UIMS通过Thrift提供RPC接口调用，pkg/thrift/client  pkg/thrift/server 为uims系统提供的客户端和服务端，client_test.go中TestInvoke作为使用client远程调用uims rpc server的单元测试，同时也提供了使用client的方法
- UIMS开启rpc服务端的方法：
```
make
./uims thrift-rpc:server
```
- 运行client单元测试
```
cd <项目跟目录>/pkg/thrift/client
go test -run=Invoke 
```
- 传参说明
- uims系统与各客户端系统约定按如下格式传参：
- 请求参数包封装成json格式，有两个基本的参数域：method_name  params 
```json
{
  "method_name": "getUserInfo",
  "params": "字典或者json编码后的字符串"
}
```
- PHP客户端使用示例详见结算系统cass-uims-rpc分支，单元测试 
``` ./vendor/bin/phpunit tests/Unit/MicoServiceAPI/EsignTest.php --filter=testGetResultByInvokeMethodViaSwoole```

- Python客户端使用示例
```
python 客户端使用相见 vzhuoserver/gen-py中client文件里的rpc_invoke方法
```

### 功能

- [x] 连接`MySQL`数据库
- [x] 测试中加载框架
- [x] 配置模块
- [] 储存模块(储存接口)
- [x] 日志模块
- [x] 多通道日志
- [x] 连接到`Redis`
- [x] 命令行工具
- [x] 数据库迁移
    - 操作流程
        1. `db/migrations` 目录下新建迁移文件, 例如 `2020_5_7_17_59_create_users_table.go`
        2. 编辑新建的文件, 实现 `Key(), Up(), Down()`
            - `Key()` 该迁移文件的唯一标识, 推荐使用文件名
            - `Up()` 执行 `migrator` 操作时会调用
            - `Down()` 执行 `migrator rollback` 操作调用
        3. 注册迁移文件到 `pkg/migrate/migrations.go -> MigrateFiles` 变量中
        4. 如果 go run main.go server 启动项目：
           - 执行迁移: `go run main.go migrator`,
           - 执行回滚: `go run main.go migrator rollback`
        5. 如果 `make run` 启动项目：
           - 执行迁移: `./uims migrator`,
           - 执行回滚: `./uims migrator rollback`
    - 手动实现的, 按步数迁移回滚操作后续补充上
    - 注意!
        - 迁移文件中不要使用 `model` 来创建表, 目的是维持迁移文件的版本性不随着 `model` 的变更而变动
- [x] 路由结构
    - 路由定义参考 `routes/api/api.go`
- [x] ORM
- [x] Swagger
    - 使用 `go build -o bin/swag vendor/github.com/swaggo/swag/cmd/swag/main.go` 生成 `bin/swag` 可执行文件
    - 使用 `./bin/swag init` 生成新的 swagger 文档
    - 通过 `http://localhost:8080/swagger/index.html` 访问 swagger 页面
- [x] 中间件
    - 参考 `routes/api/api.go` 中使用了 api 中间件组
- [x] 注册自定义表单验证规则
- [x] faker
- [] 热更新 (建议直接写测试测接口)
- [] cron 定时任务

### 相关文档

[gin](https://github.com/gin-gonic/gin)
[faker](https://github.com/bxcodec/faker)
[gorm](https://gorm.io/zh_CN/docs/)
[log](https://github.com/sirupsen/logrus)
[gin 模型验证tag文档](https://godoc.org/gopkg.in/go-playground/validator.v9)
[redis](https://github.com/go-redis/redis)
[swag](https://github.com/swaggo/swag)
[thrift](http://thrift.apache.org)
[mgo](http://labix.org/mgo)
[mgo/api](http://gopkg.in/mgo.v2)
[mgo/bson](http://gopkg.in/mgo.v2/bson)
[mgo/txn](http://gopkg.in/mgo.v2/txn)

## 开发规范

### 参考

[uber go 编码规范](https://github.com/xxjwxc/uber_go_guide_cn)

### 建议

1. commit 时必须勾选 `Perform code analysis`, `Go fmt`
2. commit 时最好勾选 `Check TODO (Show All)` 
3. 函数必须要有测试, 并且通过测试

