# UIMS系统API

#### 展示组织信息

| 接口 | 请求方式 | 说明 |
| :---: | :---: | :---: |
| /api/org/ | GET  | 展示组织信息 |

***

#### 接口详情

------

**接口地址**：/api/org/

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：此接口用于让拥有权限的角色展示组织信息。

**请求头中携带证书**：

| 参数名称 | 类型 | 是否必须 |  描述  |
| :---: | :---: | :---: | :---: |
| Authorization | string |  是   | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| organization_id | string | 是 | 组织id |

**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| id | int | 组织id |
| org_name_cn | string | 组织中文名 |
| org_name_en | string | 组织英文名 |
| org_code | string | 组织代码 |

**JSON返回示例**：

```json
{
  "code": "success",
  "sub_code": "no_auth",
  "show_msg": "",
  "debug_msg": "",
  "content": {
    "id": 1,
    "org_name_cn": "测试组织",
    "org_name_en": "test_organization",
    "org_code": "ASD123"
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
| id | 组织id |
| org_name_cn | 组织中文名 |
| org_name_en | 组织英文名 |
| org_code | 组织代码 |
