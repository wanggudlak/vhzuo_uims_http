package thriftclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wanggudlak/go-pool"
	"sync"
	"time"
	"uims/gen-go/uims_rpc_api"
	"uims/pkg/thrift/internal/hashtable"
)

const (
	jsonDataProtocol       = "json"
	binaryDataProtocol     = "binary"
	compactDataProtocol    = "compact"
	simpleJSONDataProtocol = "simplejson"
)

// Config 配置
type Config struct {
	// 服务端提供的服务，以字符串形式标明服务的位置
	ServerAPIServiceLoc string

	ServerAddr   string
	DataProtocol string
	BufferedSize int

	// socket连接超时时间
	SocketTimeout time.Duration

	// 当前配置的客户端所在的连接池标识
	poolKey int
	// 连接池中初始放入的连接实例数量
	InitialConnCountInPool int
	// 连接池最大容量
	MaxConnCountOfPool int
	// 当前客户端的连接池
	pool *pool.Pool

	Logger *log.Logger

	// 创建客户端连接实例的方法
	InstanceConnFunc pool.InstanceConn

	OnOff    bool
	Buffered bool
	Framed   bool
	Secure   bool

	// 是否启用I/O多路复用，注意：将此配置项置为true的前提是，thrift rpc server也相应的使用I/O多路复用处理器
	IsUseIOMultiplexing bool
}

// poolMap 用来存储多个客户端的连接池
var poolMap = new(sync.Map)

type Client struct {
	transport           thrift.TTransport
	rpcApiServiceClient *uims_rpc_api.UIMSRpcApiServiceClient
	config              *Config

	// 这个客户端所在的连接池
	pool *pool.Pool
	// 客户端连接实例包裹器
	connWrapper *pool.ConnWrapper
}

// Call 客户端携带请求(TRequest)向服务端发起请求，获得响应(TResponse)
func (c *Client) Call(req TRequest, resp TResponse) error {
	if c.config.IsUseIOMultiplexing {
		// 如果开启了I/O多路复用，我们并不真正关闭连接实例，只是将其重新放入连接池
		defer c.connWrapper.Close()
	} else {
		// 如果没有开启I/O多路复用，关闭连接实例，而不要放回连接池
		//c.connWrapper.MarkUnusable()
		defer c.connWrapper.Close()
	}
	r, err := c.rpcApiServiceClient.InvokeMethod(context.Background(), req.String())
	if err != nil {
		return err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return resp.Parse(b)
}

// isAllowDataProtocol 判断s是否是有效的数据协议标识
func isAllowDataProtocol(s string) bool {
	switch s {
	case jsonDataProtocol, binaryDataProtocol, compactDataProtocol, simpleJSONDataProtocol, "":
		return true
	default:
		return false
	}
}

// InitPoolKey 初始化一个连接池标识，用于标识以及寻找连接池
func (c *Config) InitPoolKey() error {
	if len(c.ServerAddr) == 0 {
		return errors.New("server address is not able empty")
	}
	key := c.ServerAddr
	c.poolKey = hashtable.Hash(key)
	return nil
}

// getDefaultInstanceClientConnFunc 获取一个默认方法，这个方法的用途是创建一个thrift rpc连接实例
func getDefaultInstanceClientConnFunc(c *Config) func() (pool.Conn, error) {
	return func() (pool.Conn, error) {
		var err error
		var conn pool.Conn
		var protocolFactory thrift.TProtocolFactory
		var transportFactory thrift.TTransportFactory
		client := new(Client)
		client.config = c
		switch c.DataProtocol {
		case compactDataProtocol:
			protocolFactory = thrift.NewTCompactProtocolFactory()
		case simpleJSONDataProtocol:
			protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
		case jsonDataProtocol:
			protocolFactory = thrift.NewTJSONProtocolFactory()
		case binaryDataProtocol, "":
			protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
		default:
			break
		}
		if nil == protocolFactory {
			return conn, fmt.Errorf("New thrift protocol factory error")
		}

		if c.Buffered {
			transportFactory = thrift.NewTBufferedTransportFactory(c.BufferedSize)
		} else {
			transportFactory = thrift.NewTTransportFactory()
		}

		if c.Framed {
			transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
		}
		if nil == transportFactory {
			return conn, fmt.Errorf("New thrift transport factory error")
		}

		if c.Secure {
			cfg := new(tls.Config)
			cfg.InsecureSkipVerify = true
			//client.transport, err = thrift.NewTSSLSocket(c.ServerAddr, cfg)
			client.transport, err = thrift.NewTSSLSocketTimeout(c.ServerAddr, cfg, c.SocketTimeout)
		} else {
			//client.transport, err = thrift.NewTSocket(c.ServerAddr)
			client.transport, err = thrift.NewTSocketTimeout(c.ServerAddr, c.SocketTimeout)
		}
		if err != nil {
			return conn, errors.Wrap(err, "New client.transport error")
		}

		client.transport, err = transportFactory.GetTransport(client.transport)
		if err != nil {
			return conn, errors.Wrap(err, "Get Transport Fail")
		}

		if err := client.transport.Open(); err != nil {
			return conn, errors.Wrap(err, "Open Transport Fail")
		}

		if c.IsUseIOMultiplexing {
			if len(c.ServerAPIServiceLoc) == 0 {
				return conn, errors.New("server api service name error")
			}
			iprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(client.transport), c.ServerAPIServiceLoc)
			oprot := thrift.NewTMultiplexedProtocol(protocolFactory.GetProtocol(client.transport), c.ServerAPIServiceLoc)
			client.rpcApiServiceClient = uims_rpc_api.NewUIMSRpcApiServiceClientProtocol(client.transport, iprot, oprot)
		} else {
			iprot := protocolFactory.GetProtocol(client.transport)
			oprot := protocolFactory.GetProtocol(client.transport)
			client.rpcApiServiceClient = uims_rpc_api.NewUIMSRpcApiServiceClient(thrift.NewTStandardClient(iprot, oprot))
		}

		return client, nil
	}
}

// Get 根据所给配置获取一个客户端实例
func Get(c *Config) (*Client, error) {
	var err error
	if !isAllowDataProtocol(c.DataProtocol) {
		return nil, errors.New("Not allow data protocol: " + c.DataProtocol)
	}

	if 0 == c.poolKey {
		if err = c.InitPoolKey(); err != nil {
			return nil, err
		}
	}

	if 0 == c.InitialConnCountInPool {
		c.InitialConnCountInPool = 5
	}
	if 0 == c.MaxConnCountOfPool {
		c.MaxConnCountOfPool = 30
	}
	if c.InitialConnCountInPool > c.MaxConnCountOfPool {
		c.InitialConnCountInPool = c.MaxConnCountOfPool
	}
	if nil == c.InstanceConnFunc {
		c.InstanceConnFunc = getDefaultInstanceClientConnFunc(c)
	}

	if v, ok := poolMap.Load(c.poolKey); ok && v != nil {
		// 我们从poolMap中获取到了关于这个客户端配置的连接池
		if ppool, ok2 := v.(*pool.Pool); !ok2 {
			ppool, err = NewClientPool(c)
			if err != nil {
				return nil, err
			}
			poolMap.Store(c.poolKey, ppool)
			return FetchOneClientInstanceFromPool(ppool)
		} else {
			return FetchOneClientInstanceFromPool(ppool)
		}
	} else {
		// 我们从poolMap中[没有]获取到当前客户端配置情形下的连接池，
		// 初始化一个连接池，并放入poolMap中
		ppool, err := NewClientPool(c)
		if err != nil {
			return nil, err
		}
		poolMap.Store(c.poolKey, ppool)
		return FetchOneClientInstanceFromPool(ppool)
	}
}

// FetchOneClientInstanceFromPool 从连接池中取出一个连接实例放入当前客户端实例中
func FetchOneClientInstanceFromPool(ppool *pool.Pool) (*Client, error) {
	if ppool != nil {
		conn, err := (*ppool).Get()
		if err != nil {
			return nil, fmt.Errorf("Get pool.Conn failed from pool: %s", err)
		}
		if cw, ok := conn.(*pool.ConnWrapper); ok {
			if pclient, ok2 := cw.Conn.(*Client); ok2 {
				pclient.connWrapper = cw
				return pclient, nil
			} else {
				return nil, errors.New("Get client failed from pool.ConnWrapper")
			}
		} else {
			return nil, errors.New("Get pool.ConnWrapper failed from pool")
		}
	} else {
		return nil, errors.New("pool.Pool is nil")
	}
}

func NewClientPool(c *Config) (*pool.Pool, error) {
	pl, err := pool.NewChannelPool(c.InitialConnCountInPool, c.MaxConnCountOfPool, c.InstanceConnFunc)
	if err != nil {
		return nil, err
	}
	return &pl, nil
}

func (c *Client) Close() error {
	if c.transport != nil && c.transport.IsOpen() {
		return c.transport.Close()
	}
	return nil
}
