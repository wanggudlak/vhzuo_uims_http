package thrift_RPC_server

import (
	"fmt"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/tcp"
	"uims/pkg/color"
	thriftserver "uims/pkg/thrift/server"
)

var CMDuimsThiftRPCServer = &command.Command{
	UsageLine: "thrift-rpc:server",
	Short:     "启动 UIMS Thrift RPC Server",
	Long: `
thrift-rpc:server 命令会启动 UIMS Thrift RPC Server。
`,
	PreRun: func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    RunningUIMSThriftRPCServer,
}

var (
	addr           string // Address to listen to
	protocol       string // Specify the protocol (binary, compact, json, simplejson)
	framed         bool   // Use framed transport
	buffered       bool   // Use buffered transport
	secure         bool   // Use tls secure transport
	ioMultiplexing bool   // 是否启用I/O多路复用
)

func init() {
	CMDuimsThiftRPCServer.Flag.StringVar(&addr, "addr", "localhost:9090", "指定服务监听地址及端口")
	CMDuimsThiftRPCServer.Flag.StringVar(&protocol, "protocol", "binary", "指定数据传输格式协议，可选值有 binary, compact, json, simplejson")
	CMDuimsThiftRPCServer.Flag.BoolVar(&framed, "framed", true, "Use framed transport")
	CMDuimsThiftRPCServer.Flag.BoolVar(&buffered, "buffered", false, "Use buffered transport")
	CMDuimsThiftRPCServer.Flag.BoolVar(&secure, "secure", true, "Use tls secure transport")
	CMDuimsThiftRPCServer.Flag.BoolVar(&ioMultiplexing, "iomul", true, "是否启用I/O多路复用，默认启用，传入false表示不启用")

	command.CMD.Register(CMDuimsThiftRPCServer)
}

// RunningUIMSThriftRPCServer 启动UIMS RPC Server
// // ./uims thrift-rpc:server -addr=0.0.0.0:9091 -buffered=false -framed=true -protocol=binary -secure=true
func RunningUIMSThriftRPCServer(cmd *command.Command, args []string) int {
	var err error
	if len(args) > 0 {
		err = cmd.Flag.Parse(args[1:])
		if err != nil {
			fmt.Println(color.Red(err.Error()))
			return 1
		}
	}
	var handle = thriftserver.NewUIMSRpcAPIHandler()
	handle.RegisterAPIwhithMap(&tcp.ThriftRPCmethodMap)
	err = thriftserver.RunningServer(handle, addr, protocol, framed, buffered, ioMultiplexing, secure)
	if err != nil {
		fmt.Println(color.Red(err.Error()))
		return 1
	}

	return 0
}
