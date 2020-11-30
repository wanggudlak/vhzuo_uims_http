# UIMS系统API

#### 0 基本信息

**正式域名1：uims.viidesk.com**

**测试域名1：test-uims.viidesk.com**

#### 1 获取已入驻UIMS系统的客户端

| 接口              | 请求方式 | 说明            |
| :----------: | :----------: | :----------: |
| /api/client/list | GET  | 展现已经入驻了UIMS系统的客户端业务系统 |

***

#### 1.1 接口详情

**接口地址**：/api/client/list

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：此接口用户UIMS系统内部，供管理员查看当前已经入驻了UIMS系统的客户端业务系统。

**请求头中携带证书**：

| 接口          | 请求方式 | 是否必须 | 描述 |
| :----------: | :----------: | :----------: | :-------: |
| Authorization | string  | 是 | 授权令牌 |

**请求参数**：

| 接口          | 请求方式 | 是否必须 | 描述 |
| :----------: | :----------: | :----------: | :-------: |
| page | 整形  | 否 | 默认值是1，取第一页的数据 |

**返回参数说明：**

|    参数名称     |  类型  | 描述        |
| :---------: | :--: | :-------- |
|  client_id  |  整形  | 客户端ID     |
| client_name | 字符串  | 客户端业务系统名称 |
|             |      |           |

**JSON返回示例**：

```json
{
  "code": "success",
  "sub_code": "",
  "show_msg": "",
  "debug_msg": "",
  "content": {
    "total": 1,
    "list": [
      {
        "client_id": 1,
        "client_name": "微桌任务系统",
    	"client_host_url": "https://marketplace.viidesk.com"
      },{
        "client_id": 2,
        "client_name": "微桌结算系统",
  		"client_host_url": "https://fuwu.skysharing.cn"
      }
  	]
  }
}
```

**返回参数及参数值说明**

|       参数值       |            描述            |
| :-------------: | :----------------------: |
|     success     |           请求成功           |
|     failed      |           请求失败           |
|    show_msg     |        显示给用户看的消息         |
|    debug_msg    | 出错时，调试信息，在测试环境下有这个返回参数的值 |
|     content     |          业务数据内容          |
|      total      |      用于分页接口中的数据总行数       |
|      list       |       用于分页接口中的list       |
|    client_id    |   客户端ID，获取客户端详情时需要传的参数   |
|   client_name   |      客户端业务系统高的中文名称       |
| client_host_url |       客户端系统的默认访问地址       |



#### 2 获取已入驻客户端系统的详情

| 接口        | 请求方式 | 说明                     |
| ----------- | -------- | ------------------------ |
| /api/client | GET      | 获取客户端系统的入驻详情 |

#### 2.1 接口详情

**接口地址**：/api/client

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：

**请求头中携带证书**：

|     参数名称      |   类型   | 是否必须 |  描述  |
| :-----------: | :----: | :--: | :--: |
| Authorization | string |  是   | 授权令牌 |

**请求参数：**

| 参数名称      | 类型   | 是否必须 | 描述      |
| :--------- | :---- | :---- | :------- |
| id | 整形   | 是    | 客户端系统ID |

**返回参数说明：**

| 参数名称              | 类型      | 描述                                       |
| :------------------ | :------- | :---------------------------------------- |
| client_name        | 字符串     | 客户端业务系统名称                                |
| client_type        | 字符串     | 客户端类型，VDK：微桌                             |
| client_flag_code   | 字符串     | 客户端业务系统标识，VDK_CASS：微桌结算系统；VDK_MP：微桌任务系统平台；VDK_CRM：微桌CRM系统；VDK_INVO：微桌代开发票系统；VDK_ESIGN：微桌电签系统；VDK_ES_SAPP：微桌电签小程序； |
| app_id             | 字符串     | 客户端系统APPID，用来唯一标识客户端系统，展示APPID，后面跟一个复制按钮 |
|                    |         |                                          |
| status             | 字符串     | 客户端业务系统使用UIMS的状态，默认N：未授权；Y：已授权；F-禁用      |
| client_host_ip     | 字符串     | 客户端当前使用的IP                               |
| client_host_url    | 字符串     | 客户端业务系统当前使用的域名，例如微桌结算系统是https://fuwu.skysharing.cn |
| in_at              | 字符串     | 入驻可以使用的开始时间点，展示格式：年-月-日 时:分:秒            |
| forget_at          | 字符串     | 在什么时间点，客户端系统不能使用UIMS，如果没有过期时间，默认是空字符串，如果有显示格式：展示格式：年-月-日 时:分:秒 |
| settings           | 列表      | 详细设置项                                    |
| setting_id         | 整形      | 设置项ID                                    |
| type               | 字符串     | 设置项类型，类型：LGN-用于登录的设置；REG-用于注册的设置；        |
| spm_full_code      | 字符串     | SPM编码                                    |
| page_template_file | 字符串     | 登录页或注册页html模板文件(.tmpl后缀)所在的位置(UIMS系统的根目录的相对路径，并以appid作为子目录名)，因为涉及多阶段登录，采用json存储，{"a":"/downloads/app_id/login_a.tmpl", "b": "/sownloads/app_id/login_b.tmpl"} |
| form_fields        | JSON字符串 | 表单域属性数据                                  |

**JSON响应示例：**

```json
{
  "code": "success",
  "sub_code": "",
  "show_msg": "",
  "debug_msg": "",
  "content": {
    "client_name": "客户端业务系统名称",
    "client_type": "VDK",
    "client_flag_code": "VDK_CASS",
    "app_id": "34718923748172384",
    "status": "Y",
    "client_host_ip": "39.106.127.99",
    "client_host_url": "https://fuwu.skysharing.cn",
    "in_at": "2020-04-27 00:00:00",
    "forget_at": "",
    "settings": [
      {
        "setting_id": 1,
        "type": "LGN",
        "spm_full_code": "1024.34718923748172384.100.101",
        "page_template_file": [
          "a": "/downloads/app_id/login_a.tmpl"
        ],
        "form_fields": "{\"a\": [{\"attr_id\": \"account\", \"attr_cn\": \"账号\"},{\"attr_id\": \"passwd\", \"attr_cn\": \"密码\"}, {\"attr_id\": \"sms_code\", \"attr_cn\": \"验证码\", \"type\": \"phone_sms\"}]}"
      }
    ]
  }
}
```

