



# UIMS系统API

#### 用户基础信息相关接口

|      接口      | 请求方式 |             说明             |
| :------------: | :------: | :--------------------------: |
| /api/user/list |   GET    | 获取用户列表(仅展示基础数据) |
|   /api/user    |   POST   |           用户注册           |
|   /api/user    |   PUT    |   修改用户状态(封禁或解禁)   |

***

#### 接口详情

------

**接口地址:/api/users/list

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：获取用户列表(仅展示基础数据)

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 |    描述    |
| :------: | :--: | :------: | :--------: |
|   page   | Int  |    否    | 默认值是1  |
| pagesize | Int  |    否    | 默认值是10 |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": [
        {
            "id": 1,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "13545046834",
            "wechat": "",
            "wechat_id": "omICa09K0G8akO-aUL6enGqb2ZtY",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:50+08:00",
            "updated_at": "2020-05-07T17:49:50+08:00"
        },
        {
            "id": 2,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "18627704982",
            "wechat": "",
            "wechat_id": "omICa0xi06j0YoFoiJtMPlCKV3tw",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:51+08:00",
            "updated_at": "2020-05-07T17:49:51+08:00"
        },
        {
            "id": 3,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "15071127549",
            "wechat": "",
            "wechat_id": "",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:52+08:00",
            "updated_at": "2020-05-07T17:49:52+08:00"
        },
        {
            "id": 4,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "15827263638",
            "wechat": "",
            "wechat_id": "",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:52+08:00",
            "updated_at": "2020-05-07T17:49:52+08:00"
        },
        {
            "id": 5,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "18627981216",
            "wechat": "",
            "wechat_id": "",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:53+08:00",
            "updated_at": "2020-05-07T17:49:53+08:00"
        },
        {
            "id": 6,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "15827174344",
            "wechat": "",
            "wechat_id": "omICa09yWw9Ar4Yr16TV4Vb_y1rc",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:54+08:00",
            "updated_at": "2020-05-07T17:49:54+08:00"
        },
        {
            "id": 7,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "15572532156",
            "wechat": "",
            "wechat_id": "",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:54+08:00",
            "updated_at": "2020-05-07T17:49:54+08:00"
        },
        {
            "id": 8,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "18971420040",
            "wechat": "",
            "wechat_id": "omICa0zvOc1qp52U7Ij7P9FHzLIo",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:55+08:00",
            "updated_at": "2020-05-07T17:49:55+08:00"
        },
        {
            "id": 9,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "",
            "wechat": "",
            "wechat_id": "",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:56+08:00",
            "updated_at": "2020-05-07T17:49:56+08:00"
        },
        {
            "id": 10,
            "user_type": "VDK",
            "account": "",
            "user_code": "",
            "na_code": "+86",
            "phone": "19986920641",
            "wechat": "",
            "wechat_id": "",
            "email": "",
            "encrypt_type": 1,
            "status": "Y",
            "created_at": "2020-05-07T17:49:56+08:00",
            "updated_at": "2020-05-07T17:49:56+08:00"
        }
    ]
}
```

**接口地址:/api/users

**请求方式**：Post

**请求和响应数据格式**：JSON

**接口备注**：注册用户

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

|   参数名称   |  类型  | 是否必须 |          描述          |
| :----------: | :----: | :------: | :--------------------: |
|     page     |  Int   |    是    |         手机号         |
|    email     | string |    否    |          邮箱          |
|   account    | string |    否    |       自定义账号       |
|    passwd    | string |    是    |          密码          |
| encrypt_type |  int   |    是    | 加密方式(1:微桌2:结算) |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "id": 2772,
        "user_type": "VDK",
        "account": "18571818278",
        "user_code": "1262936625084633088",
        "na_code": "+86",
        "phone": "18571818278",
        "wechat": "",
        "wechat_id": "",
        "email": "aizedi23@sina.com",
        "encrypt_type": 2,
        "status": "Y",
        "created_at": "2020-05-20T10:42:10.30524+08:00",
        "updated_at": "2020-05-20T10:42:10.30524+08:00"
    }
}
```

**接口地址** : /api/resources

**请求方式**：PUT

**请求和响应数据格式**：JSON

**接口备注**：修改用户状态

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述         |
| -------- | ---- | -------- | ------------ |
| user_id  | int  | 是       | 用户id       |
| type     | int  | 是       | 1:解禁2:禁用 |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": null
}
```



_______

**资源点相关接口返回参数说明:**

|   参数名称   |                       参数说明                       |
| :----------: | :--------------------------------------------------: |
|      id      |                       资源的id                       |
|   res_code   |                       资源编码                       |
|   account    |                 和前端约定的资源编码                 |
|   res_type   |          资源类型，A：逻辑资源；B：实体资源          |
| res_sub_type | 源子类型，AP：页面；AC：菜单；AM：按钮；AD：数据资源 |
|    wechat    |                    资源的英文名称                    |
|  wechat_id   |                    资源的中文名称                    |
| encrypt_type |                  资源的后端路由URI                   |
|    status    |                      N正常Y禁用                      |
|    isdel     |        是否删除: 默认N：未软删除；Y：已软删除        |
|  created_at   |                       创建时间                       |



**通用返回参数及参数值说明**

|     参数值      |                       描述                       |
| :-------------: | :----------------------------------------------: |
|     success     |                     请求成功                     |
|     failed      |                     请求失败                     |
|    show_msg     |                显示给用户看的消息                |
|    debug_msg    | 出错时，调试信息，在测试环境下有这个返回参数的值 |
|     content     |                   业务数据内容                   |
|   total_pages   |              用于分页接口中的总页数              |
|      list       |               用于分页接口中的list               |
|    client_id    |      客户端ID，获取客户端详情时需要传的参数      |
|   client_name   |            客户端业务系统高的中文名称            |
| client_host_url |             客户端系统的默认访问地址             |