package service

import (
	"encoding/json"
	"errors"
	"fmt"
	thriftserver "uims/pkg/thrift/server"
)

func GetUserService() UserService {
	return UserService{}
}

// 角色
func GetRoleService() RoleService {
	return RoleService{}
}

//资源点
func GetResourceService() ResourceService {
	return ResourceService{}
}

//资源组
func GetResGroupService() ResGroupService {
	return ResGroupService{}
}

//角色关联资源组
func GetRoleResMapService() RoleResMapService {
	return RoleResMapService{}
}

//客户端
func GetClientService() ClientService {
	return ClientService{}
}

//客户端设置信息
func GetClientSettingService() ClientSettingService {
	return ClientSettingService{}
}

//组织
func GetOrgService() OrgService {
	return OrgService{}
}

//组织
func GetUserRoleService() UserRoleService {
	return UserRoleService{}
}

// Thrift client
func GetThriftClientServer() ThriftClientServer {
	return ThriftClientServer{}
}

// Thrift server
func DemoRPCbizHandleOld(v interface{}) (string, error) {
	responseBizBody := map[string]interface{}{}
	bizParamBody := map[string]interface{}{}
	switch v.(type) {
	case string:
		fmt.Println("v is string:", v)
		err := json.Unmarshal([]byte(v.(string)), &bizParamBody)
		if err != nil {
			return "", errors.New(fmt.Sprintf("业务请求参数解码失败：<%s>", err.Error()))
		}
		// ... 调用具体业务逻辑方法
		fmt.Println(bizParamBody)

	case map[string]interface{}:
		// ... 调用具体业务逻辑方法

		fmt.Println("v is map[string]interface{}", v)

	default:
		return "", errors.New("业务请求参数格式错误")
	}

	responseBizBody["biz_content"] = "aleijuzixiaodou"
	responseBizBody["biz_status"] = "lovecoder"
	responseBizBodyBytes, err := json.Marshal(responseBizBody)

	return string(responseBizBodyBytes), err
}

type DemoReq struct {
	Test string `json:"test" binding:"required"`
}

func DemoRPCbizHandleNew(c *thriftserver.Context) {
	var req DemoReq
	if err := c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	// .. 调用具体业务逻辑方法
	c.Response.Success("业务响应参数", "业务处理完成")
	return
}
