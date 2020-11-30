package service

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"uims/conf"
	"uims/pkg/glog"
	thriftclient "uims/pkg/thrift/client"
)

type ThriftClientServer struct{}

// client info
type client struct {
	ID             uint   `json:"id"`
	ClientType     string `json:"client_type"`
	ClientFlagCode string `json:"client_flag_code"`
	ClientName     string `json:"client_name"`
}

// UIMS is client to subsystem is server send data scheduler
func (t ThriftClientServer) ClientInvoke(clientId int, apiMethodName string, items interface{}) Response {
	resp := Response{}
	err := t.Invoke(Request{
		BRequest: thriftclient.BRequest{
			MethodName: apiMethodName,
			Params:     items,
		},
	}, &resp, clientId)
	if err != nil {
		resp.Status = "failed"
		resp.Msg = err.Error()
	}
	return resp
}

type Request struct {
	thriftclient.BRequest
}

type Response struct {
	thriftclient.BResponse
}

func (t ThriftClientServer) InvokeMP(req Request, resp thriftclient.TResponse) {
	req.Params = map[string]interface{}{
		"items": req.Params,
	}
	resp.SetStatus("failed")

	config := conf.ThriftClients["mp"]

	// Thrift Client on-off is true
	if config.OnOff == false {
		resp.SetStatus("thrift client on-off is off")
		resp.SetMsg("(MP) thrift client on-off is off")
		return
	}
	// 循环请求10次
	for i := 0; i < 10; i++ {
		cli, err := thriftclient.Get(&config)
		if err != nil {
			glog.Channel("thrift").WithError(err).WithFields(log.Fields{
				"request": req,
			}).Error("获取 thrift client 实例失败")
			resp.SetMsg(err.Error())
			return
		}
		err = cli.Call(req, resp)
		if err != nil {
			glog.Channel("thrift").WithError(err).WithFields(log.Fields{
				"request": req,
			}).Error("请求 RPC 失败")
			resp.SetMsg(err.Error())
			continue
		}
		if resp.CallOK() {
			return
		}
	}
	return
}

func (ThriftClientServer) Invoke(req Request, resp *Response, clientId int) error {
	resp.Status = "failed"
	resp.Msg = ""
	clientInfo, err := GetClientService().GetClientByID(clientId)

	if err != nil {
		return err
	}

	var config thriftclient.Config
	switch clientInfo.ClientType {
	case "CASS": // 微桌结算系统
		config = conf.ThriftClients["cass"]
	case "VDK": // 微桌任务系统
		config = conf.ThriftClients["mp"]
	default:
		return errors.New("service platform does not exist")
	}

	// Thrift Client on-off is true
	if config.OnOff == false {
		resp.SetStatus("thrift client on-off is off")
		resp.SetMsg(clientInfo.ClientName + "（" + clientInfo.ClientType + "）thrift client on-off is off")
		return nil
	}
	req.Params = map[string]interface{}{
		"items": req.Params,
		"client": &client{
			ID:             clientInfo.ID,
			ClientType:     clientInfo.ClientType,
			ClientFlagCode: clientInfo.ClientFlagCode,
			ClientName:     clientInfo.ClientName,
		},
	}
	// 循环请求10次
	for i := 0; i < 10; i++ {
		cli, err := thriftclient.Get(&config)
		if err != nil {
			return err
		}
		err = cli.Call(req, resp)
		if err != nil {
			glog.Channel("thrift").WithError(err).WithFields(log.Fields{
				"request": req,
			}).Error("请求 RPC 失败")
			resp.Msg = err.Error()
			continue
		}
		if resp.GetStatus() == "success" {
			return nil
		}
	}

	return errors.New(req.MethodName + " data synchronization failed")
}
