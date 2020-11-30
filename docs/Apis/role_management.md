## UIMS系统API

#### 角色管理
| 接口  | 请求方式 | 说明 |
| :--- | :---: | :--- |
|[/api/roles/list ](#RoleListGet)| GET | 获取客户端角色列表or客户端用户角色列表 |
|[/api/roles ](#RoleGet)|  GET | 获取角色信息 |
|[/api/roles ](#RolePost)| POST | 添加角色 |
|[/api/roles ](#RolePUt)| PUT | 更新角色 |
|[/api/roles ](#RoleDelete)| DELETE | 删除角色 |
|[/api/role/user/list](#RoleUserListGet)| GET | 获取角色用户列表 |
|[/api/role/user](#RoleUserPost)| Post | 角色添加用户 |
|[/api/role/user](#RoleUserDelete)| Delete | 角色移除用户 |

***

#### 接口详情
***

**<span name="RoleListGet">获取角色列表</span>** 

**接口地址**：/api/roles/list

**请求方式**：GET

**返回格式**：JSON

**接口备注**：此接口供管理员查看当前客户端角色列表，客户端用户的角色列表。

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| page | int | 否 | 页码数 |
| pagesize | int | 否 | 一页的数量 |
| client_id | int | 是 | 客户端id |
| user_id | int | 是 | user_id |
| org_id | int | 否 | 组织id |

**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| error_code | int | 状态码 |
| page | int | 否 | 页码数 |
| size | int | 否 | 一页的数量,默认10 |
| tota | int | 角色总数|
| roles | list | 角色列表 |

**JSON返回示例**：

```json
{
    "code": "success",
    "sub_code": "",
    "show_msg": "",
    "debug_msg": "",
    "content": {
    "page":1,
    "size":10
        "total": 5,
        "roles":[
       		 {"id": "",  "client_id": "", "name_cn": "", "name_en": "",  "org_id": ""}
       		 		]
           		 			}  
}
```



***

#### 接口详情
***

**<span name="RolePost">获取角色信息</span>** 

**接口地址**：/api/role

**请求方式**：GET

**返回格式**：JSON

**接口备注**：获取角色信息

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| id | int | 是 | 角色id |



**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| id | int | 否 | 页码数 |
| client_id | int | int | 客户端id |
| name_cn | int | 角色名称|
| name_en | string | 角色英文名称 |
| org_id | int | 组织id |

**JSON返回示例**：

```json
{
    "code": "success",
    "sub_code": "",
    "show_msg": "",
    "debug_msg": "",
    "content": {
            "id": "",
            "client_id": "",
            "name_cn": "",
            "name_en": "",
            "org_id": "",
            }  
}
```

***

#### 接口详情
***

**<span name="RolePost">添加角色</span>** 

**接口地址**：/api/role

**请求方式**：Post

**返回格式**：JSON

**接口备注**：添加角色

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| client_id | int | 否 | 客户端id |
| org_id | int | 否 | 组织id |
| name_cn | int | 是 | 角色名称 |
| name_en | int | 是 | 角色英文名称 |


**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| id | int | 否 | 页码数 |
| client_id | int | int | 客户端id |
| name_cn | int | 角色名称|
| name_en | string | 角色英文名称 |
| org_id | int | 组织id |

**JSON返回示例**：

```json
{
    "code": "success",
    "sub_code": "",
    "show_msg": "",
    "debug_msg": "",
    "content": {
            "id": "",
            "client_id": "",
            "name_cn": "",
            "name_en": "",
            "org_id": "",
            }  
}
```



***

#### 接口详情
***

**<span name="RolePut">更新角色</span>** 

**接口地址**：/api/role

**请求方式**：Put

**返回格式**：JSON

**接口备注**：更新角色

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| role_id | int | 否 | 角色id |
| name_cn | int | 是 | 角色名称 |
| name_en | int | 是 | 角色英文名称 |
| org_id | int  | 否 | 组织id |

**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| id | int | 否 | 页码数 |
| client_id | int | int | 客户端id |
| name_cn | int | 角色名称|
| name_en | string | 角色英文名称 |
| org_id | int | 组织id |

**JSON返回示例**：

```json
{
    "code": "success",
    "sub_code": "",
    "show_msg": "",
    "debug_msg": "",
    "content": {
            "id": "",
            "client_id": "",
            "name_cn": "",
            "name_en": "",
            "org_id": "",
            }  
}
```

***

#### 接口详情
***

**<span name="RoleDelete">删除角色</span>** 

**接口地址**：/api/role

**请求方式**：Delete

**返回格式**：JSON

**接口备注**：删除角色

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| role_id | int | 否 | 角色id |


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
    "content": {
    	"id":1
    }  
}
```

***

#### 接口详情
***

**<span name="RoleUserListGet">获取角色下的所有用户</span>** 

**接口地址**：/api/role/user/list

**请求方式**：GET

**返回格式**：JSON

**接口备注**：用于获取角色下面的所有用户

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| id   | int | 否 | 角色id |


**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| account | string | 状态码描述信息 |
| id | string | 状态码描述信息 |


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
        "user": [
                {"id": 1, "account": "" }
               		 ]
                			}  
}
```


**

#### 接口详情
***

**<span name="RoleUserPost">角色添加用户</span>** 

**接口地址**：/api/role/user

**请求方式**：POST

**返回格式**：JSON

**接口备注**：用于角色添加用户

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| id   | int | 否 | 角色id |
| user_id   | int | 是 | 用户id |


**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| id | string | 状态码描述信息 |


**JSON返回示例**：

```json
{
    "code": "success",
    "sub_code": "",
    "show_msg": "",
    "debug_msg": "",
    
    "content": {
            "id": 1
           	 }  
}
```

***

#### 接口详情
***

**<span name="RoleUserDelete">角色移除用户</span>** 

**接口地址**：/api/role/user

**请求方式**：DELETE

**返回格式**：JSON

**接口备注**：用于角色移除用户

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Authorization | string | 是 | 授权令牌 |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| id   | int | 否 | 角色id |
| user_id   | int | 是 | 用户id |


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
    "content": {
    	"id": 1
    	}
}
```

