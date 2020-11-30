package thriftserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"uims/gen-go/uims_rpc_api"
	"uims/pkg/color"
	"uims/pkg/thrift/common"
)

const BUFFER_SIZE = 8192

// RunningServer
func RunningServer(apiHandler *UIMSRpcAPIHandler, addr, protocol string, framed, buffered, ioMultiplexing, secure bool) error {
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case common.COMPACT:
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case common.SIMPLE_JSON:
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case common.JSON:
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case common.BINARY:
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		return errors.New("Invalid data protocol")
	}

	var transportFactory thrift.TTransportFactory
	if buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(BUFFER_SIZE)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	var err error
	var transport thrift.TServerTransport

	if secure {
		cfg := new(tls.Config)
		if cert, err := LoadX509KeyPair(common.KEY_DIR+"/server.crt", common.KEY_DIR+"/server.key", common.PRIVATE_KEY_PASSWD); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}

	if err != nil {
		return err
	}

	if ioMultiplexing {
		processor := thrift.NewTMultiplexedProcessor()
		processor.RegisterProcessor("UIMSRpcApiService", uims_rpc_api.NewUIMSRpcApiServiceProcessor(apiHandler))
		server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
		fmt.Println(color.Green(fmt.Sprintf("Starting the thrift rpc server... on %s", addr)))
		return server.Serve()
	} else {
		processor := uims_rpc_api.NewUIMSRpcApiServiceProcessor(apiHandler)
		server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
		fmt.Println(color.Green(fmt.Sprintf("Starting the thrift rpc server... on %s", addr)))
		return server.Serve()
	}
}
