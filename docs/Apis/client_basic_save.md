# UIMS系统API

#### 保存入驻的基本信息数据

| 接口 | 请求方式 | 说明 |
| :---: | :---: | :---: |
| /api/client | POST  | 保存入驻的基本信息数据 |

***

#### 接口详情

------

**接口地址**：/api/client

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：此接口用于让拥有权限的角色进行保存基本信息数据。

**请求头中携带证书**：

| 参数名称 | 类型 | 是否必须 |  描述  |
| :---: | :---: | :---: | :---: |
| Authorization | string |  是   | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| client_type | string | 是 | 客户端类型，VDK：微桌 |
| client_flag_code | string | 是 | 客户端业务系统标识，VDK_CASS：微桌结算系统等； |
| client_spm1_code | string | 是 | SPM编码中的第一部分，微桌内部系统用1024；外部系统用2048 |
| client_spm2_code | string | 是 | SPM编码中的第二部分 |
| client_name | string | 是 | 客户端业务系统名称 |
| client_host_ip | json | 是 | 客户端当前使用的IP，多个用json字符串保存 |
| client_host_url | json | 是 | 客户端业务系统当前使用的域名，多个用json字符串保存 |
| client_pub_key_path | string | 是 | 客户端业务系统的RSA公钥key文件路径 |
| in_at | datetime | 否 | 入驻可以使用的开始时间点，默认为当前时间 |
| forget_at | datetime | 否 | 在什么时间点，客户端系统不能使用UIMS |

**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| id | int | 客户端ID |

**JSON返回示例**：

```json
{
  "code": "success",
  "sub_code": "no_auth",
  "show_msg": "",
  "debug_msg": "",
  "content": {
    "id": 1
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
| id | 客户端ID |
