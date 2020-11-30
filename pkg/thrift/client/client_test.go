package thriftclient

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var config = Config{
	ServerAddr: "192.168.50.215:9092",
	//ServerAddr:   "127.0.0.1:9091",
	DataProtocol: "binary",
	BufferedSize: 8192,
	Buffered:     false,
	Framed:       true,
	Secure:       false,

	ServerAPIServiceLoc:    "UIMSRpcApiService",
	InitialConnCountInPool: 5,
	MaxConnCountOfPool:     30,
	SocketTimeout:          time.Second * 10,
	IsUseIOMultiplexing:    true,
	Logger:                 log.StandardLogger(),
}

func TestGet(t *testing.T) {
	client, err := Get(&config)
	assert.Nil(t, err)

	body := BRequest{
		MethodName: "get_uims_wx_qr",
		Params: map[string]interface{}{
			"items": map[string]string{},
		},
	}
	for i := 0; i < 10; i++ {
		// 连续发送消息
		r, err := client.rpcApiServiceClient.InvokeMethod(context.Background(), body.String())
		t.Log(r)
		assert.Nil(t, err)
	}
}

type Resp struct {
	BResponse
}

func TestCli_Call(t *testing.T) {
	for i := 0; i < 10000; i++ {
		client, err := Get(&config)
		assert.Nil(t, err)
		resp := Resp{}
		err = client.Call(BRequest{
			MethodName: "test",
			Params: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
		}, &resp)
		t.Logf("%+v", resp)
		assert.Nil(t, err)
	}
}

func TestBody_String(t *testing.T) {
	body := BRequest{
		MethodName: "test",
		Params: map[string]string{
			"test": "test",
		},
	}
	t.Log(body.String())
	assert.NotEmpty(t, body)

	body = BRequest{
		MethodName: "test",
		Params:     nil,
	}
	t.Log(body.String())
	assert.NotEmpty(t, body)
}

func TestConfig_InitPoolKey(t *testing.T) {
	err := (&config).InitPoolKey()
	if err != nil {
		t.Errorf("Init pool key error %s", err)
	}
	t.Log(config.poolKey)
}
