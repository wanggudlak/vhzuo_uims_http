# UIMS系统API

#### 资源策略组相关接口

|           接口           | 请求方式 |          说明          |
| :----------------------: | :------: | :--------------------: |
| /api/resource/group/list |   GET    |    获取用户资源组数    |
|   /api/resource/group    |   GET    |     获取资源组详情     |
|   /api/resource/group    |   POST   |     添加资源组数据     |
|   /api/resource/group    |   PUT    |     修改资源组数据     |
|   /api/resource/group    |  DELETE  | 删除资源点数据(假删除) |

***

#### 接口详情

------

**接口地址:/api/resource/group/list

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：获取用户资源组数据列表

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

| 参数名称  | 类型 | 是否必须 |    描述    |
| :-------: | :--: | :------: | :--------: |
|   page    | Int  |    否    | 默认值是1  |
| pagesize  | Int  |    否    | 默认值是10 |
| client_id | int  |    是    | 客户端id       |
|  user_id  | int  |    否    | 用户id    |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "count": 3,
        "source_group_list": [
            {
                "client_id": 1,
                "created_at": 1589439135,
                "id": 5,
                "isdel": "N",
                "org_id": 1,
                "res_group_cn": "测试",
                "res_group_code": "1231231",
                "res_group_en": "test",
                "res_group_type": "DEFAULT",
                "res_of_curr": {
                    "resource_ids": [
                        1,
                        2,
                        3,
                        4
                    ]
                },
                "updated_at": 1589439135
            },
            {
                "client_id": 1,
                "created_at": 1589353577,
                "id": 3,
                "isdel": "N",
                "org_id": 0,
                "res_group_cn": "运营经理",
                "res_group_code": "",
                "res_group_en": "",
                "res_group_type": "DEFAULT",
                "res_of_curr": {
                    "resource_ids": [
                        289,
                        290,
                        291
                    ]
                },
                "updated_at": 1589353582
            },
            {
                "client_id": 1,
                "created_at": 1588850239,
                "id": 2,
                "isdel": "N",
                "org_id": 1,
                "res_group_cn": "会员管理333",
                "res_group_code": "22212121",
                "res_group_en": "vip_manage333",
                "res_group_type": "AC",
                "res_of_curr": {
                    "resource_ids": [
                        290,
                        29,
                        292
                    ]
                },
                "updated_at": 1588850239
            }
        ]
    }

}
```



**接口地址:/api/resource/group

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：获取资源组详情

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 |   描述   |
| :------: | :--: | :------: | :------: |
|    id    | Int  |    否    | 资源组id |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "resource_group_info": {
            "client_id": 1,
            "created_at": 1588850239,
            "id": 2,
            "isdel": "N",
            "org_id": 1,
            "res_group_cn": "会员管理333",
            "res_group_code": "2222",
            "res_group_en": "vip_manage333",
            "res_group_type": "AC",
            "res_of_curr": {
                "resource_ids": [
                    290,
                    29,
                    292
                ]
            },
            "updated_at": 1588850239
        },
        "resource_list": [
            {
                "client_id": 1,
                "created_at": 1588848546,
                "id": 29,
                "org_id": 0,
                "res_code": "VHypFYBn",
                "res_data_location": null,
                "res_endp_route": "",
                "res_front_code": "kRgjA5nF",
                "res_name_cn": "删除帐户交易明细",
                "res_name_en": "delete_trade",
                "res_type": "A"
            },
            {
                "client_id": 1,
                "created_at": 1588848546,
                "id": 290,
                "org_id": 0,
                "res_code": "qzKv61OL",
                "res_data_location": null,
                "res_endp_route": "",
                "res_front_code": "4O0JG8nw",
                "res_name_cn": "导出回访记录",
                "res_name_en": "export_visit_record",
                "res_type": "A"
            },
            {
                "client_id": 1,
                "created_at": 1588848546,
                "id": 292,
                "org_id": 0,
                "res_code": "PLxngyR6",
                "res_data_location": null,
                "res_endp_route": "",
                "res_front_code": "frpDuFjL",
                "res_name_cn": "更新作品审核",
                "res_name_en": "update_works_review",
                "res_type": "A"
            }
        ]
    }
}
```

**接口地址** : /api/resource/group

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：添加资源点数据

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

|     参数名称      |                          参数说明                           |
| :---------------: | :---------------------------------------------------------: |
|     client_id     |                      客户端业务系统id                       |
|      org_id       |                        客户端组织id                         |
|  res_group_code   |                         资源组编码                          |
| resources_id_list |                   资源点列表:默认[]空数组                   |
|  res_group_type   | 权限策略组类型：DEFAULT-默认策略组；SELF-自定义配置的策略组 |
|   res_group_en    |                      资源组的英文名称                       |
|   res_group_cn    |                      资源组的中文名称                       |



**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "client_id": 2,
        "created_at": 1588850239,
        "id": 2,
        "isdel": "N",
        "org_id": 1,
        "res_group_cn": "会员管理333",
        "res_group_code": "22212121",
        "res_group_en": "vip_manage333",
        "res_group_type": "AC",
        "res_of_curr": {
            "resource_ids": [
                290,
                29,
                292
            ]
        },
        "updated_at": 1588850239
    }
}
```

**接口地址** : /api/resource/group

**请求方式**：PUT

**请求和响应数据格式**：JSON

**接口备注**：修改资源点数据

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

|     参数名称      |                          参数说明                           |
| :---------------: | :---------------------------------------------------------: |
|  res_group_code   |                         资源组编码                          |
| resources_id_list |                       资源点列表数组                        |
|  res_group_type   | 权限策略组类型：DEFAULT-默认策略组；SELF-自定义配置的策略组 |
|   res_group_en    |                      资源组的英文名称                       |
|   res_group_cn    |                      资源组的中文名称                       |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "client_id": 2,
        "created_at": 1588850239,
        "id": 2,
        "isdel": "N",
        "org_id": 1,
        "res_group_cn": "会员管理333",
        "res_group_code": "22212121",
        "res_group_en": "vip_manage333",
        "res_group_type": "AC",
        "res_of_curr": {
            "resource_ids": [
                290,
                29,
                292
            ]
        },
        "updated_at": 1588850239
    }
}
```

**接口地址:/api/resource/group

**请求方式**：DELETE

**请求和响应数据格式**：JSON

**接口备注**：删除资源组

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 |   描述   |
| :------: | :--: | :------: | :------: |
|    id    | int  |    是    | 资源组id |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": "nil"
}
```



_______

**资源组相关接口返回参数说明:**

|     参数名称      |                          参数说明                           |
| :---------------: | :---------------------------------------------------------: |
|        id         |                         资源组的id                          |
|  res_group_code   |                         资源组编码                          |
| resources_id_list |                   所包含的资源点数据列表                    |
|  res_group_type   | 权限策略组类型：DEFAULT-默认策略组；SELF-自定义配置的策略组 |
|   res_group_en    |                      资源组的英文名称                       |
|      org_id       |                           组织id                            |
|   res_group_cn    |                      资源组的中文名称                       |
|       isdel       |           是否删除: 默认N：未软删除；Y：已软删除            |
|     created_at     |                          创建时间                           |



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