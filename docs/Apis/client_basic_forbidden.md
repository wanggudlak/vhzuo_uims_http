# UIMS系统API

#### 禁用入驻的基本信息数据

| 接口 | 请求方式 | 说明 |
| :---: | :---: | :---: |
| /api/client/status | POST  | 修改客户端状态 Y 授权  F 禁用 |

***

#### 接口详情

------

**接口地址**：/api/client/status

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：此接口用于让拥有权限的角色禁用入驻的基本信息数据。

**请求头中携带证书**：

| 参数名称 | 类型 | 是否必须 |  描述  |
| :---: | :---: | :---: | :---: |
| Authorization | string |  是   | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| id | int | 是 | 客户端id |
| status | string | 是 |状态 |

**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| is_success | bool | 是否成功，true-是，false-否 |

**JSON返回示例**：

```json
{
  "code": "success",
  "sub_code": "no_auth",
  "show_msg": "",
  "debug_msg": "",
  "content": {
    "is_success": true
  }
}
```

**返回参数及参数值说明**

| 参数值 | 描述 |
| :---: | :---: |
| success | 请求成功 |
| failed | 请求失败 |
| show_msg | 显示给用户看的消息 |
| debug_msg | 出错时，调试信息，在测试环境下有这个返回参数的值 |
| content | 业务数据内容 |
| is_success | 是否成功，true-是，false-否 |
