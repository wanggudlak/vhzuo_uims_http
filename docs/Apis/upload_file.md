# UIMS系统API

#### 上传文件

| 接口 | 请求方式 | 说明 |
| :---: | :---: | :---: |
| /api/client/file/upload | POST  | 上传文件 |

***

#### 接口详情

------

**接口地址**：/api/client/file/upload

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：此接口用于上传文件。

**请求头中携带证书**：

| 参数名称 | 类型 | 是否必须 |  描述  |
| :---: | :---: | :---: | :---: |
| Authorization | string |  是   | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| file_type | string | 是 | 上传文件类型：1-商户公钥，2-html模版 |
| file_data | file | 是 | 文件内容 |

**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| file_path | string | 文件相对路径 |

**JSON返回示例**：

```json
{
  "code": "success",
  "sub_code": "no_auth",
  "show_msg": "",
  "debug_msg": "",
  "content": {
    "file_path": "/upload/publicKey/123124234(app_id名称)/210134123.pem"
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
| file_path | 文件相对路径 |
