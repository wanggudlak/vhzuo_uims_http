# UIMS系统API

#### 资源点相关数据接口

|        接口         | 请求方式 |          说明          |
| :-----------------: | :------: | :--------------------: |
| /api/resource/list |   GET    | 获取用户资源点数据列表 |
|   /api/resource    |   GET    |     获取资源点详情     |
|   /api/resource    |   POST   |     添加资源点数据     |
|   /api/resource    |   PUT    |     修改资源点数据     |
|   /api/resource    |  DELETE  | 删除资源点数据(假删除) |

***

#### 接口详情

------

**接口地址:/api/resource/list

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：获取用户资源点数据列表

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

| 参数名称  | 类型 | 是否必须 |    描述    |
| :-------: | :--: | :------: | :--------: |
|   page    | Int  |    否    | 默认值是1  |
| pagesize  | Int  |    否    | 默认值是10 |
| client_id | int  |    是    |  客户端id  |
|  user_id  | int  |    否    |   用户id   |
|  role_id  | int  |    否    |   角色id   |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "data_list": [
            {
                "id": 291,
                "client_id": 1,
                "org_id": 0,
                "res_code": "KWPYEwST",
                "res_front_code": "GrwlpZdL",
                "res_type": "A",
                "res_sub_type": "AD",
                "res_name_en": "insert_works_review",
                "res_name_cn": "添加作品审核",
                "res_endp_route": "",
                "res_data_location": null,
                "is_del": "N",
                "created_at": 1588848546,
                "updated_at": 1588848546
            },
            {
                "id": 292,
                "client_id": 1,
                "org_id": 0,
                "res_code": "PLxngyR6",
                "res_front_code": "frpDuFjL",
                "res_type": "A",
                "res_sub_type": "AD",
                "res_name_en": "update_works_review",
                "res_name_cn": "更新作品审核",
                "res_endp_route": "",
                "res_data_location": null,
                "is_del": "N",
                "created_at": 1588848546,
                "updated_at": 1588848546
            },
            {
                "id": 293,
                "client_id": 1,
                "org_id": 0,
                "res_code": "yn0UlzMv",
                "res_front_code": "sWk8gqN5",
                "res_type": "A",
                "res_sub_type": "AD",
                "res_name_en": "select_works_review",
                "res_name_cn": "查看作品审核",
                "res_endp_route": "",
                "res_data_location": null,
                "is_del": "N",
                "created_at": 1588848546,
                "updated_at": 1588848546
            },
            {
                "id": 294,
                "client_id": 1,
                "org_id": 0,
                "res_code": "xW9cL0bm",
                "res_front_code": "Co32GN8n",
                "res_type": "A",
                "res_sub_type": "AD",
                "res_name_en": "delete_works_review",
                "res_name_cn": "删除作品审核",
                "res_endp_route": "",
                "res_data_location": null,
                "is_del": "N",
                "created_at": 1588848546,
                "updated_at": 1588848546
            },
            {
                "id": 295,
                "client_id": 1,
                "org_id": 0,
                "res_code": "F7JnRNck",
                "res_front_code": "JLrw3ZfP",
                "res_type": "A",
                "res_sub_type": "AD",
                "res_name_en": "export_works_review",
                "res_name_cn": "导出作品审核",
                "res_endp_route": "",
                "res_data_location": null,
                "is_del": "N",
                "created_at": 1588848546,
                "updated_at": 1588848546
            },
            {
                "id": 296,
                "client_id": 1,
                "org_id": 1,
                "res_code": "1",
                "res_front_code": "1",
                "res_type": "A",
                "res_sub_type": "AC",
                "res_name_en": "vip_manage",
                "res_name_cn": "会员管理",
                "res_endp_route": "/vip/manage1/",
                "res_data_location": {
                    "database": "",
                    "table": "",
                    "status": ""
                },
                "is_del": "N",
                "created_at": 1588927213,
                "updated_at": 1588927213
            },
            {
                "id": 297,
                "client_id": 1,
                "org_id": 1,
                "res_code": "bqqhvalh5s5nief0rep0",
                "res_front_code": "bqqhvalh5s5nief0repg",
                "res_type": "A",
                "res_sub_type": "AC",
                "res_name_en": "vip_manage",
                "res_name_cn": "会员管理",
                "res_endp_route": "/vip/manage1/",
                "res_data_location": {
                    "database": "",
                    "table": "",
                    "status": ""
                },
                "is_del": "N",
                "created_at": 1588928427,
                "updated_at": 1588928427
            },
            {
                "id": 298,
                "client_id": 1,
                "org_id": 1,
                "res_code": "bqqis9dh5s5gbmcbnf00",
                "res_front_code": "bqqis9dh5s5gbmcbnf0g",
                "res_type": "A",
                "res_sub_type": "AC",
                "res_name_en": "vip_manage",
                "res_name_cn": "会员管理",
                "res_endp_route": "/vip/manage1/",
                "res_data_location": {
                    "database": "",
                    "table": "",
                    "status": ""
                },
                "is_del": "N",
                "created_at": 1588932133,
                "updated_at": 1588932133
            },
            {
                "id": 299,
                "client_id": 1,
                "org_id": 1,
                "res_code": "bqqj0clh5s5inpa85m0g",
                "res_front_code": "bqqj0clh5s5inpa85m10",
                "res_type": "A",
                "res_sub_type": "AC",
                "res_name_en": "vip_manage",
                "res_name_cn": "会员管理",
                "res_endp_route": "/vip/manage1/",
                "res_data_location": {
                    "database": "",
                    "table": "user",
                    "status": "normal"
                },
                "is_del": "N",
                "created_at": 1588932658,
                "updated_at": 1588932658
            },
            {
                "id": 300,
                "client_id": 1,
                "org_id": 1,
                "res_code": "bqqj3k5h5s5j2ff3vrn0",
                "res_front_code": "bqqj3k5h5s5j2ff3vrng",
                "res_type": "A",
                "res_sub_type": "AC",
                "res_name_en": "vip_manage",
                "res_name_cn": "会员管理",
                "res_endp_route": "/vip/manage1/",
                "res_data_location": {
                    "database": "uims",
                    "table": "user",
                    "status": "normal"
                },
                "is_del": "N",
                "created_at": 1588933072,
                "updated_at": 1588933072
            }
        ],
        "total_num": 306
    }
}
```

**返回参数说明:**

|     参数名称      |                           参数说明                           |
| :---------------: | :----------------------------------------------------------: |
|        id         |                           资源的id                           |
|     res_code      |                           资源编码                           |
|  res_front_code   |                     和前端约定的资源编码                     |
|     res_type      |              资源类型，A：逻辑资源；B：实体资源              |
|   res_sub_type    |     源子类型，AP：页面；AC：菜单；AM：按钮；AD：数据资源     |
|    res_name_en    |                        资源的英文名称                        |
|    res_name_cn    |                        资源的中文名称                        |
|  res_endp_route   |                      资源的后端路由URI                       |
| res_data_location | 资源所在的位置，主要用于数据权限，json存储，包含以下属性：客户端id、数据库名、表名等 |
|       isdel       |            是否删除: 默认N：未软删除；Y：已软删除            |
|     created_at     |                           创建时间                           |



**接口地址:/api/resource

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：获取资源点详情

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 |   描述    |
| :------: | :--: | :------: | :-------: |
|   id   | Int  |    否    | 默认值是1 |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "id": 300,
        "client_id": 1,
        "org_id": 1,
        "res_code": "bqqj3k5h5s5j2ff3vrn0",
        "res_front_code": "bqqj3k5h5s5j2ff3vrng",
        "res_type": "A",
        "res_sub_type": "AC",
        "res_name_en": "vip_manage",
        "res_name_cn": "会员管理",
        "res_endp_route": "/vip/manage1/",
        "res_data_location": {
            "database": "uims",
            "table": "user",
            "status": "normal"
        },
        "is_del": "N",
        "created_at": 1588933072,
        "updated_at": 1588933072
    }
}
```

**接口地址** : /api/resource

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：添加资源点数据

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：
|     参数名称      |类型 |是否必传|                          参数说明                           |
| :---------------: | ----|-----|:----------------------------------------------------------: |
|
|     client_id     |   int  |是|                  客户端业务系统id                       |
|      org_id       |   int  |是|                    客户端组织id                         |
|     res_code      | string   |是|                        资源编码                           |
|  res_front_code   |       string|是|              和前端约定的资源编码                     |
|     res_type      |     string |是|        资源类型，A：逻辑资源；B：实体资源              |
|   res_sub_type    |  string|是|   源子类型，AP：页面；AC：菜单；AM：按钮；AD：数据资源     |
|    res_name_en    |   string   |是|                  资源的英文名称                        |
|    res_name_cn    |     string|是|                   资源的中文名称                        |
|  res_endp_route   |        string|是|              资源的后端路由URI                       |
| res_data_location | string|是|资源所在的位置，主要用于数据权限，json存储，包含以下属性：客户端id、数据库名、表名等 |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "id": 300,
        "client_id": 1,
        "org_id": 1,
        "res_code": "bqqj3k5h5s5j2ff3vrn0",
        "res_front_code": "bqqj3k5h5s5j2ff3vrng",
        "res_type": "A",
        "res_sub_type": "AC",
        "res_name_en": "vip_manage",
        "res_name_cn": "会员管理",
        "res_endp_route": "/vip/manage1/",
        "res_data_location": {
            "database": "uims",
            "table": "user",
            "status": "normal"
        },
        "is_del": "N",
        "created_at": 1588933072,
        "updated_at": 1588933072
    }
}
```

**接口地址** : /api/resource

**请求方式**：PUT

**请求和响应数据格式**：JSON

**接口备注**：修改资源点数据

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

|     参数名称      |类型 |是否必传|                          参数说明                           |
| :---------------: | ----|-----|:----------------------------------------------------------: |
|        id         | int    |是|                      资源点id                           |
|     res_code      | string   |是|                        资源编码                           |
|  res_front_code   |       string|是|              和前端约定的资源编码                     |
|     res_type      |     string |是|        资源类型，A：逻辑资源；B：实体资源              |
|   res_sub_type    |  string|是|   源子类型，AP：页面；AC：菜单；AM：按钮；AD：数据资源     |
|    res_name_en    |   string   |是|                  资源的英文名称                        |
|    res_name_cn    |     string|是|                   资源的中文名称                        |
|  res_endp_route   |        string|是|              资源的后端路由URI                       |
| res_data_location | string|是|资源所在的位置，主要用于数据权限，json存储，包含以下属性：客户端id、数据库名、表名等 |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "id": 300,
        "client_id": 1,
        "org_id": 1,
        "res_code": "bqqj3k5h5s5j2ff3vrn0",
        "res_front_code": "bqqj3k5h5s5j2ff3vrng",
        "res_type": "A",
        "res_sub_type": "AC",
        "res_name_en": "vip_manage",
        "res_name_cn": "会员管理",
        "res_endp_route": "/vip/manage1/",
        "res_data_location": {
            "database": "uims",
            "table": "user",
            "status": "normal"
        },
        "is_del": "N",
        "created_at": 1588933072,
        "updated_at": 1588933072
    }
}
```

**接口地址:/api/resource

**请求方式**：DELETE

**请求和响应数据格式**：JSON

**接口备注**：删除资源点

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 |   描述   |
| :------: | :--: | :------: | :------: |
|    id    | int  |    是    | 资源点id |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": "nil"
}
```



_______

**资源点相关接口返回参数说明:**

|     参数名称      |                           参数说明                           |
| :---------------: | :----------------------------------------------------------: |
|        id         |                           资源的id                           |
|     res_code      |                           资源编码                           |
|  res_front_code   |                     和前端约定的资源编码                     |
|     res_type      |              资源类型，A：逻辑资源；B：实体资源              |
|   res_sub_type    |     源子类型，AP：页面；AC：菜单；AM：按钮；AD：数据资源     |
|    res_name_en    |                        资源的英文名称                        |
|    res_name_cn    |                        资源的中文名称                        |
|  res_endp_route   |                      资源的后端路由URI                       |
| res_data_location | 资源所在的位置，主要用于数据权限，json存储，包含以下属性：客户端id、数据库名、表名等 |
|       isdel       |            是否删除: 默认N：未软删除；Y：已软删除            |
|     created_at     |                           创建时间                           |



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