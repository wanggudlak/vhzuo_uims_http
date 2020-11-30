# UIMS系统API

#### 上传文件

| 接口 | 请求方式 | 说明 |
| :---: | :---: | :---: |
| /api/background_login | POST  | 后台登陆 |

***

#### 接口详情
------

**接口地址**：/api/background_login

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：此接口用于上传文件。

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| account | string |  是   | 账号 |
| passwd | string |  是   | 密码 |

**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :--- | :---: | :---: |
| code | integer | code码 |
| message | string | 返回文案 |
| body | object | 响应数据 |
| └token | string | 登陆token |
| └user | object | 登陆用户信息数据 |
| └└account | string | 登陆账号 |
| └└created_at | string | 创建时间 |
| └└email | string | 邮箱 |
| └└encrypt_type | integer | 加密类型 |
| └└id | integer | 主键id |
| └└na_code | string | 国家代码 |
| └└open_id | string | 用户唯一ID标示 |
| └└open_id | string | 用户唯一ID标示 |
| └└phone | string | 手机号 |
| └└status | string | 状态 |
| └└updated_at | string | 更新时间 |
| └└user_code | string | 用户编码 |
| └└user_type | string | 用户类型 |

**JSON返回示例**：

```json
{
  "code": 0,
  "message": "success",
  "body": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoidWltc19zdXBlcl9hZG1pbiIsImV4cCI6MTU5NDc0NDI1MiwiaXNzIjoiVUlNUy1CQUNLIiwibmJmIjoxNTk0NzIxNjUyfQ.AQWHRcS7d1orXJBFitNHQ7GdWTlOJkElcRePNYHtY_s",
      "user": {
        "account": "uims_super_admin",
        "created_at": "2020-07-14T16:16:41+08:00",
        "email": "vzhuo@vzhuo.com",
        "encrypt_type": 0,
        "id": 1,
        "na_code": "+86",
        "open_id": "",
        "phone": "17852000001",
        "status": "Y",
        "updated_at": "2020-07-14T16:16:41+08:00",
        "user_code": "1282952138284077056",
        "user_type": "CASS"
      }
  }
}
```

**返回参数及参数值说明**

| 参数值 | 描述 |
| :---: | :---: |
| success | 请求成功 |
| failed | 请求失败 |
| code | 返回状态码 |
| body | 业务数据内容 |
| token | token |
| user | 用户信息 |
| account | 账号 |
| created_at | 创建时间 |
| email | 邮箱 |
| id | 主键id |
| na_code | 国家编码 |
| open_id | 用户唯一标示 |
| phone | 手机号 |
| status | 状态 |
| user_code | 用户编码 |
| user_type | 用户类型 |
