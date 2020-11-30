



# UIMS系统API

#### 用户资料相关接口

|      接口      | 请求方式 |             说明             |
| :------------: | :------: | :--------------------------: |
| /api/user/info |   GET    | 获取用户列表(仅展示基础数据) |
| /api/user/info |   POST   |           用户注册           |
| /api/user/info |   PUT    |   修改用户状态(封禁或解禁)   |

***

#### 接口详情

**请求方式**：Post

**请求和响应数据格式**：JSON

**接口备注**：创建用户资料

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

|         参数名称         |  类型  | 是否必须 |       描述       |
| :----------------------: | :----: | :------: | :--------------: |
|         birthday         | string |    是    |   "2019-02-12"   |
|      header_img_url      | string |    是    |   头像地址url    |
| identity_card_emblem_img | string |    否    | 身份证国徽面地址 |
|     identity_card_no     | string |    否    |     身份证号     |
| identity_card_person_img | string |    否    |     身份账号     |
|      landline_phone      | string |    否    |      座机号      |
|         nickname         | string |    是    |       昵称       |
|           sex            | string |    是    |     M:男W:女     |
|        taxer_type        | string |    否    |    纳税人类型    |
|         taxer_no         | string |    否    |     纳税编号     |
|         user_id          |  int   |    是    |      用户id      |
|         name_cn          | string |    否    |      中文名      |
|         name_en          | string |    否    |      英文名      |
|      name_cn_alias       | string |    否    |     中文别名     |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "birthday": "2019-02-02",
        "created_at": "2020-05-20T14:59:59.583804+08:00",
        "header_img_url": "www.baidu.com",
        "id": 2126,
        "identity_card_emblem_img": "./upload/img_path",
        "identity_card_no": "22220022020",
        "identity_card_person_img": "./upload/img_path",
        "is_identify": "",
        "isdel": "",
        "landline_phone": "",
        "na_code": "",
        "name_abbr_py": "a",
        "name_cn": "艾艾艾",
        "name_cn_alias": "",
        "name_en": "",
        "name_full_py": "aiaiai",
        "nickname": "aimoly",
        "sex": "",
        "taxer_no": "",
        "taxer_type": "",
        "updated_at": "2020-05-20T14:59:59.583805+08:00",
        "user_code": "",
        "user_id": 33,
        "user_type": ""
    }
}
```

**接口地址** : /api/user/info

**请求方式**：PUT

**请求和响应数据格式**：JSON

**接口备注**：修改用户资料

**请求头中携带证书**：

|   参数名称    |  类型  | 是否必须 |   描述   |
| :-----------: | :----: | :------: | :------: |
| Authorization | string |    是    | 授权令牌 |

**请求参数**：

|         参数名称         |  类型  | 是否必须 |       描述       |
| :----------------------: | :----: | :------: | :--------------: |
|         birthday         | string |    是    |   "2019-02-12"   |
|      header_img_url      | string |    是    |   头像地址url    |
| identity_card_emblem_img | string |    否    | 身份证国徽面地址 |
|     identity_card_no     | string |    否    |     身份证号     |
| identity_card_person_img | string |    否    |     身份账号     |
|      landline_phone      | string |    否    |      座机号      |
|         nickname         | string |    是    |       昵称       |
|           sex            | string |    是    |     M:男W:女     |
|        taxer_type        | string |    否    |    纳税人类型    |
|         taxer_no         | string |    否    |     纳税编号     |
|         user_id          |  int   |    是    |      用户id      |
|         name_cn          | string |    否    |      中文名      |
|         name_en          | string |    否    |      英文名      |
|      name_cn_alias       | string |    否    |     中文别名     |

**JSON返回示例**：

```json
{
    "code": 0,
    "message": "success",
    "body": {
        "birthday": "2019-02-02",
        "created_at": "2020-05-20T14:59:59.583804+08:00",
        "header_img_url": "www.baidu.com",
        "id": 2126,
        "identity_card_emblem_img": "./upload/img_path",
        "identity_card_no": "22220022020",
        "identity_card_person_img": "./upload/img_path",
        "is_identify": "",
        "isdel": "",
        "landline_phone": "",
        "na_code": "",
        "name_abbr_py": "a",
        "name_cn": "艾艾艾",
        "name_cn_alias": "",
        "name_en": "",
        "name_full_py": "aiaiai",
        "nickname": "aimoly",
        "sex": "",
        "taxer_no": "",
        "taxer_type": "",
        "updated_at": "2020-05-20T14:59:59.583805+08:00",
        "user_code": "",
        "user_id": 33,
        "user_type": ""
    }
}
```



_______

**资源点相关接口返回参数说明:**

| 参数名称                 | 描述             |
| ------------------------ | ---------------- |
| birthday                 | 生日             |
| created_at                | 创建时间         |
| updated_at                | 更新时间         |
| header_img_url           | 头像             |
| id                       | 主键id           |
| identity_card_emblem_img | 身份证国徽面     |
| identity_card_no         | 身份证号         |
| identity_card_person_img | 身份证人像面     |
| is_identify              | 是否认证Y是N否   |
| landline_phone           | 座机号           |
| na_code                  | +86              |
| name_abbr_py             | 中文名拼音首字母 |
| name_cn                  | 中文名称         |
| name_cn_alias            | 中文别名         |
| name_en                  | 英文名称         |
| name_full_py             | 中文拼音全拼     |
| nickname                 | 昵称             |
| sex                      | 性别: M男W女     |
| taxer_no                 | 纳税人编号       |
| taxer_type               | 纳税人类型       |
| user_code                | 用户唯一编码     |
| user_id                  | 用户id           |

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