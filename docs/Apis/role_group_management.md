### UIMS系统API

#### 角色资源组管理
| 接口  | 请求方式 | 说明 |
|[/api/role/res/group](#RoleResGET)| GET | 角色资源组信息获取 |
|[/api/role/res/group](#RoleRespost)| POST | 角色添加资源组 |
|[/api/role/res/group](#RoleResDelete)| DELETE | 角色移除资源组 |

***

#### 接口详情
***

**<span name="RoleResGET">角色资源组信息获取</span>** 

**接口地址**：/api/role/res/group

**请求方式**：GET

**返回格式**：JSON

**接口备注**：角色资源组信息获取

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| page | int | 否 | 页码数 |
| pagesize | int | 否 | 一页的数量 |
| id   | int | 否 | 角色id |


**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| res_group_code | string | 资源组代码 |
| res_group_cn | string | 资源组中文名称 |
| res_group_en | string | 资源组英文名称 |
| id | string | 资源组id |



**JSON返回示例**：

```json
{
    "code": "success",
    "sub_code": "",
    "show_msg": "",
    "debug_msg": "",
    "content": {
        "total": 5,
        "count":5,
        "group":[
        				{id": "", "res_group_code":"","res_group_cn":"", "res_group_en":"" }
        				]
}
```

***

#### 接口详情
***

**<span name="RoleRespost">角色添加资源组</span>** 

**接口地址**：/api/role/res/group

**请求方式**：POST

**返回格式**：JSON

**接口备注**：用于角色添加资源组

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| id   | int | 否 | 角色id |
| res_id   | int | 是 | 资源组id |


**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| res_group_code | string | 资源组代码 |
| res_group_cn | string | 资源组中文名称 |
| res_group_en | string | 资源组英文名称 |
| id | string | 资源组id |



**JSON返回示例**：

```json
{
    "code": "success",
    "sub_code": "",
    "show_msg": "",
    "debug_msg": "",
    "content": {
    	"id": "",
        "res_group_code":"",
        "res_group_cn":"",   
        "res_group_en":"",
               			 }  
}
```

***

#### 接口详情
***

**<span name="RoleResDelete">角色移除资源组</span>** 

**接口地址**：/api/role/res/group

**请求方式**：Delete

**返回格式**：JSON

**接口备注**：角色移除资源组

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| id   | int | 否 | 角色id |
| res_id  | int | 是 | 资源组id |


**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |


**JSON返回示例**：

```json
{
    "code": "success",
    "sub_code": "",
    "show_msg": "",
    "debug_msg": "",
    "content": {}  
}
```