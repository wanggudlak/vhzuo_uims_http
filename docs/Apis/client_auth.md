# Client 鉴权

## 名词

`客户端`: 指使用 UIMS 系统的服务, 例如代付系统, 任务系统
`用户`: 指客户端下的用户
`客户端页面`: 服务自身提供的web页面服务
`UIMS页面`: 通常指登录页或者注册页, 这些页面是在 UIMS 域名下管理的
`code`: 第一步登录后获取到的 hash 字符串, 用于客户端从 UIMS 中获取用户 Token
`state`: 客户端跳转到登录页面带上的 state 业务场景字符串, 登录第3步跳转回客户端页面时, 将原样返回

## 登录流程 

1. 用户 uims login page 输入: account, password
2. login page 发起请求到 uims 进行第一步登录
    - account
    - password
    - state
    - redirect_uri
    - form_id
    - response_type(固定`code`)
3. uims 响应 login 结果, 包含 redirect_uri 301, 跳转后到客户端页面中
    - state
    - code
4. 客户端页面拿到 code, 通过客户端请求 UIMS , 获取到用户的 Token 信息
    - code 在 UIMS 中鉴定时, 会判断 state 是否符合, code 是否有效
    - state
5. Token 返回给前端保存, 登录完成

## 接口列表

| 接口               | 请求方式 | 说明                    |
| :--------------- | :--: | :-------------------- |
| ADMIN API              |      | UIMS 后台使用 API                      |
| /api/admin/users/{user_uuid}/status | PUT  | 1.1 冻结解冻用户  |
| WEB API |   | UIMS web 调用  |
| /web/users | POST  | 2.1 用户 form 注册 |
| /web/users/login | POST  | 2.2 用户 form 登录, 返回 code + state |
| CLIENT API           |        | 客户端调用 |
| /api/client/token/code/auth | GET  | 3.1 客户端通过 code + state 换取用户的登录 token |
| /api/client/token/user | GET  | 3.2 客户端通过 token 获取用户的信息(鉴别token是否有效/获取用户身份) |
| /api/client/token/refresh | GET  | 3.3 客户端通过 refresh_token 获取新的 token |
| /api/client/users | POST  | 3.4 用户 api 注册 |
| /api/client/users/{user_uuid}/detail | GET  | 3.5 查询用户信息 |
| /api/client/users/{user_uuid} | PUT  | 3.6 完善用户信息  |
| /api/client/users/{user_uuid}/password | PUT  | 3.7 用户修改密码  |

***

## 接口详情

### ADMIN (uims 后台)

#### 1.1 冻结解冻用户

------

**接口地址**：/api/users/{user_uuid}/status

**请求方式**：PUT

**请求和响应数据格式**：JSON

**接口备注**：

冻结或者解冻用户

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   user_uuid   |   string   |   是   | 用户唯一uuid     |
|   status   |   string   |   是   | 设置的状态值 Y 正常 N 冻结 |

**返回参数说明**

无业务返回参数

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|    | |  |  |

**JSON返回示例**：

```json
{
  "token_type": "Bearer",
  "expires_in": 1588227438,
  "access_token": "kjk;643asdfga345sadt234",
  "refresh_token": "346dfthe5yhdfgyer65hr6u",
  "user_uuid": "3456345sdfg3gs34gw43g"
}
```

### WEB

#### 2.1 用户注册

------

**接口地址**：/web/users

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：

uims web 注册流程

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   account   |   string   |   否   | 用户登录流程中获取到的 code     |
|   phone   |   string   |   否   | 用户登录流程中获取到的 state     |
|   password   |   string   |   否  | 登录使用密码     |
|   verify_code   |   string   |   否   | 登录使用验证码     |
|   state   |   string   |   是   | 客户端自定义, 将会在跳转会客户端页面时带回去    |
|   redirect_uri   |   string   |   是   | 登录成功后301跳转地址     |
|   form_id   |   string   |   是   | 使用的表单id, 渲染页面时提供     |

**返回参数说明**

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   is_success   | int(1)     | 1 注册成功, 0 注册失败     | 是 |
|   redirect_uri   | string     | 网页需要跳转的地址, 与请求时一致 | 是 |

**JSON返回示例**：

```json
{
  "is_success": 1,
  "redirect_uri": "http://fuwu.skysharing.cn"
}
```

#### 2.2 用户登录

------

**接口地址**：/web/users/login

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：

登录流程第一步, 用户提交form表单

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   account   |   string   |   否   | 用户登录流程中获取到的 code     |
|   phone   |   string   |   否   | 用户登录流程中获取到的 state     |
|   password   |   string   |   否  | 登录使用密码     |
|   verify_code   |   string   |   否   | 登录使用验证码     |
|   state   |   string   |   是   | 客户端自定义, 将会在跳转会客户端页面时带回去    |
|   redirect_uri   |   string   |   是   | 登录成功后301跳转地址     |
|   form_id   |   string   |   是   | 使用的表单id, 渲染页面时提供     |
|   response_type   |   string   |   是   | 固定使用 `code`     |

**返回参数说明**

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   redirect_uri   | string     | 网页需要跳转的地址, 后面会带上 code, state 参数 | 是 |

**JSON返回示例**：

```json
{
  "redirect_uri": "http://fuwu.skysharing.cn?code=42dnsfgh24346erg&state=login"
}
```

### CLIENT

#### 3.1 code 换取 token

------

**接口地址**：/api/token/code/auth

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：

登录流程的第四步调用接口

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   code   |   string   |   是   | 用户登录流程中获取到的 code     |
|   state   |   string   |   是   | 用户登录流程中获取到的 state     |

**返回参数说明**

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   token_type   | string     | Token 类型, 目前为 Bearer 或 Mac, 大小写不明感     | 是 |
|   expires_in   | int(11)     | 到期时间, 时间戳 | 是 |
|   access_token   | string   | Token 字符串     | 是 |
|   refresh_token   | string   | 更新令牌, 用来获取下一次访问令牌 | 否 |
|   user_uuid   | string   | 用户的唯一标识 | 是 |

**JSON返回示例**：

```json
{
  "token_type": "Bearer",
  "expires_in": 1588227438,
  "access_token": "kjk;643asdfga345sadt234",
  "refresh_token": "346dfthe5yhdfgyer65hr6u",
  "user_uuid": "3456345sdfg3gs34gw43g"
}
```

#### 3.2 token 获取用户身份

---

**接口地址**：/api/token/user

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：

当token跨服务使用时, 新的服务可以使用这个接口来鉴别token的有效性, 以及区分用户.

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   token   |   string   |   是   | 用户使用的token     |
|   token_type   |   string   |   是   | 用户使用的token type     |

**返回参数说明**

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   token_type   | string     | Token 类型, 目前为 Bearer 或 Mac, 大小写不明感     | 是 |
|   expires_in   | int(11)     | 到期时间, 时间戳 | 是 |
|   access_token   | string   | Token 字符串, 与传入的一致     | 是 |
|   refresh_token   | string   | 更新令牌, 用来获取下一次访问令牌 | 否 |
|   user_uuid   | string   | 用户的唯一标识 | 是 |

**JSON返回示例**：

```json
{
  "token_type": "Bearer",
  "expires_in": 1588227438,
  "access_token": "kjk;643asdfga345sadt234",
  "refresh_token": "346dfthe5yhdfgyer65hr6u",
  "user_uuid": "3456345sdfg3gs34gw43g"
}
```

#### 3.3 通过 refresh_token 更新 token

---

**接口地址**：/api/token/refresh

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：

token 快过期时, 客户端可以使用 refresh_token 获取新 token

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   token   |   string   |   是   | 用户使用的token     |
|   token_type   |   string   |   是   | 用户使用的token type     |

**返回参数说明**

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   token_type   | string     | Token 类型, 目前为 Bearer 或 Mac, 大小写不明感     | 是 |
|   expires_in   | int(11)     | 到期时间, 时间戳 | 是 |
|   access_token   | string   | Token 字符串, 与传入的一致     | 是 |
|   refresh_token   | string   | 更新令牌, 用来获取下一次访问令牌 | 否 |
|   user_uuid   | string   | 用户的唯一标识 | 是 |

**JSON返回示例**：

```json
{
  "token_type": "Bearer",
  "expires_in": 1588227438,
  "access_token": "kjk;643asdfga345sadt234",
  "refresh_token": "346dfthe5yhdfgyer65hr6u",
  "user_uuid": "3456345sdfg3gs34gw43g"
}
```

#### 3.4 API 注册用户

---

**接口地址**：/api/users

**请求方式**：POST

**请求和响应数据格式**：JSON

**接口备注**：

客户端通过 API 注册用户

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   account   |   string   |   否   | 用户使用的token     |
|   phone   |   string   |   否   | 用户使用的token type     |
|   na_code   |   string   |   否   | 国家代码: 中国 +86, 默认 +86     |
|   passwd   |   string   |   是   | 密码     |
|   status   |   string   |   否   | 账号状态: Y：正常；N：已冻结. 默认 N     |

**返回参数说明**

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   user_uuid   | string     | 注册成功用户唯一标识 | 是 |

**JSON返回示例**：

```json
{
  "user_uuid": "3456345sdfg3gs34gw43g"
}
```

#### 3.5 获取用户详细信息

---

**接口地址**：/api/users/{user_uuid}/detail

**请求方式**：GET

**请求和响应数据格式**：JSON

**接口备注**：

根据 user_uuid 查询用户信息

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   user_uuid   |   string   |   是   | 用户唯一 user_uuid |

**返回参数说明**

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   user_uuid   | string   | 用户的唯一标识 | 是 |
|   account   | string     | 登录使用账号   | 否 |
|   user_code   | int(11)     | 用户代码, 在一个组织下唯一 | 是 |
|   na_code   | string   | 国家代码: 中国 +86     | 否 |
|   phone   | string   | 更新令牌, 用来获取下一次访问令牌 | 否 |
|   wechat   | string   | 微信账号 | 否 |
|   email   | string   | 邮箱 | 否 |
|   status   | string   | 用户状态, Y：正常；N：已冻结 |  |
|   is_identity   | string   | 是否已经实名认证, N：没有；Y：已经实名认证 |  |
|   name_en   | string   | 英文名 | 否 |
|   name_cn   | string   | 中文名 | 否 |
|   name_cn_alias   | string   | 中文别名 | 否 |
|   name_abbr_py   | string   | 姓名拼音首字母 | 否 |
|   name_full_py   | string   | 姓名拼音全拼 | 否 |
|   identity_card_no   | string   | 认证的身份证号 | 否 |
|   landline_phone   | string   | 座机号码 | 否 |
|   sex   | string   | 性别, M：男；F：女 | 否 |
|   age   | string   | 出生年月日 | 否 |
|   nickname   | string   | 昵称 | 否 |
|   taxer_type   | string   | 纳税人类型，A：一般纳税人 | 否 |
|   taxer_no   | string   | 纳税人识别号 | 否 |
|   header_img_full_url   | string   | 头像完整访问链接 | 否 |

**JSON返回示例**：

```json
{
    "user_uuid":  "66835277896744975JLIPCMA44OMJDKL",
    "account": "accountt",
    "user_code":  "123456",
    "na_code": "+86",
    "phone":  "13517210000",
    "wechat": "张三",
    "email":  "zhangsan@gmail.com",
    "status": "Y",
    "is_identity":  "Y",
    "name_en": "Zhang San",
    "name_cn":  "张三",
    "name_cn_alias": "三哥",
    "name_abbr_py":  "zs",
    "name_full_py": "zhang san",
    "identity_card_no":  "420222202005011000",
    "landline_phone": "02715155558",
    "sex":  "M",
    "age": "2020-05-01",
    "nickname":  "张三丰",
    "taxer_type": "A",
    "taxer_no":  "0978927894ABDFS",
    "header_img_full_url": "https://uims.skysharing.cn/avater/2020/05/01/453565346435.png"
}
```

#### 3.6 更新用户详细信息

---

**接口地址**：/api/users/{user_uuid}

**请求方式**：PUT

**请求和响应数据格式**：JSON

**接口备注**：

更新用户信息资料

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   user_uuid   |   string   |   是   | 用户唯一 user_uuid |

**返回参数说明**

路由参数: user_uuid 用户唯一标识

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   na_code   | string   | 国家代码: 中国 +86     | 否 |
|   is_identity   | string   | 是否已经实名认证, N：没有；Y：已经实名认证 |  |
|   name_en   | string   | 英文名 | 否 |
|   name_cn   | string   | 中文名 | 否 |
|   name_cn_alias   | string   | 中文别名 | 否 |
|   name_abbr_py   | string   | 姓名拼音首字母 | 否 |
|   name_full_py   | string   | 姓名拼音全拼 | 否 |
|   identity_card_no   | string   | 认证的身份证号 | 否 |
|   landline_phone   | string   | 座机号码 | 否 |
|   sex   | string   | 性别, M：男；F：女 | 否 |
|   age   | string   | 出生年月日 | 否 |
|   nickname   | string   | 昵称 | 否 |
|   taxer_type   | string   | 纳税人类型，A：一般纳税人 | 否 |
|   taxer_no   | string   | 纳税人识别号 | 否 |
|   header_img_base64   | string   | 头像图片 base64 编码 | 否 |

**JSON返回示例**：

```json
{
    "na_code": "+86",
    "is_identity":  "Y",
    "name_en": "Zhang San",
    "name_cn":  "张三",
    "name_cn_alias": "三哥",
    "name_abbr_py":  "zs",
    "name_full_py": "zhang san",
    "identity_card_no":  "420222202005011000",
    "landline_phone": "02715155558",
    "sex":  "M",
    "age": "2020-05-01",
    "nickname":  "张三丰",
    "taxer_type": "A",
    "taxer_no":  "0978927894ABDFS",
    "header_img_base64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAR0AAADSCAIAAAD...=="
}
```

#### 3.7 修改用户密码

------

**接口地址**：/api/client/users/{user_uuid}/password

**请求方式**：PUT

**请求和响应数据格式**：JSON

**接口备注**：

修改用户密码

**请求头**：

**请求参数**：

| 参数名称 |  类型  | 是否必须 |  描述  |
| :--: | :--: | :--: | :--: |
|   new_passwd   |   string   |   是   | 用户新密码 |
|   verify_code   |   string   |   是   | 验证码, 通过 通用接口发送 |



**返回参数说明**

| 参数名称 |  类型  |  描述  | 是否必填 |
| :--: | :--: | :--: | :--: |
|   is_success   |   int(1)   |  是否操作成功, 1成功, 0失败 | 是 |

**JSON返回示例**：

```json
{
  "new_passwd": "新密码",
  "verify_code": "543633"
}
```
