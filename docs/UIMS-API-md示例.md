#### UIMS系统API

| 接口  | 请求方式 | 说明 |
| :--- | :---: | :--- |
|[/api/v2/rpo/apply](#RpoApplyGet)| GET | RPO项目应聘列表 |

***

#### 接口详情
***

**<span name="RpoApplyGet">RPO项目应聘列表</span>** 

**接口地址**：/api/v2/rpo/apply

**请求方式**：GET

**返回格式**：JSON

**接口备注**：需求方查询RPO项目应聘列表

**请求头**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| Cookie | string | 是 | 已登录用户的session |

**请求参数**：

| 参数名称 | 类型 | 是否必须 | 描述 |
| :---: | :---: | :---: | :---: |
| pagenum | int | 否 | 页码数 |
| pagesize | int | 否 | 一页的数量 |
| job_id | string | 否 | RPO项目UUID |

**返回参数说明**

| 参数名称 | 类型 | 描述 |
| :---: | :---: | :---: |
| msg | string | 状态码描述信息 |
| error_code | int | 状态码 |
| data | list | 应聘列表详情数据 |
| count | int | 总数目数 |
| pagenum | int | 页码数 |
| statistics | json | 面试者/发送offer等数据统计 |

**JSON返回示例**：

```json
{
    "msg": "ok",
    "data": [
        {"user": 
            {
                "id": "",
                "name": "",
                "real_name": "",
                "gender": "",
                "is_student": "",
                "age": "",
                "avatar": "",
                "highest_edu": ""
             },
        "proposal": {
                "id": "", 
                "status": "", 
                "ptype": "",
                "recommand_user": [
                    {
                        "id": "324sdfer23e",    /** 推荐者的UUID **/
                        "name": "",             /** 推荐者的昵称 **/
                        "real_name": "",        /** 推荐者的真实名称 **/
                        "company_name": "",     /** 推荐者的所属企业名称 **/
                        "is_choice_rpo": true,  /** 是否是当前推荐者 **/
                    }
                ]
            },
        "interview": {
                "one": {
                    "interview_id": "",
                    "status": "",
                    "interview_type": "",
                    "interview_at": "",
                    "message": "",
                    "created_at": "",
                    "interview_addr: ""
                    },
                "two": {
                    "interview_id": "",
                    "status": "",
                    "interview_type": "",
                    "interview_at": "",
                    "message": "",
                    "created_at": "",
                    "interview_addr: ""
                    }
                },
        "offer": {
            "offer_id": "",
            "status": "",
            "budget": "",
            "credential": ""
        },
        "entry": {
            "entry_id": "",
            "entry_at": "",
            "follow": "",
            "entry_status": "",
            "information": ""
        },
        "cur_action": ""    /** 当前状态 **/
        }
    ],
    "error_code": 0,
    "count": 0,
    "pagenum": 1,
    "statistics": {
        "interview": 0,     /** 邀请面试中的人 **/
        "review": 0,        /** 邀请二面中的人 **/
        "apply": 0,         /** 应聘中的人 **/
        "offer": 0,         /** 发送过offer的人 **/
        "entry": 0,         /** 第三方入职的人 **/
        "sure_entry": 0,    /** 确认入职的人 **/
        "amount": 0,        /** 支付服务费用 **/
    }
}
```

**返回参数值 cur_action 说明**

| 参数值 | 描述 |
| :---: | :---: |
| proposal_refuse | pass应聘 |
| proposal_out | RPO服务方未被选中 |
| interview_invite | 面试邀请 |
| interview_pass | pass面试 |
| review | 邀请二面 |
| review_pass | pass二面 |
| send_offer | 发送offer |
| job_close | 项目关闭 |
| wait_entry | 已上传入职信息 |
| ensure | 质保期内 |
| expire | 已过质保期 |
| half | 支付50% |
| paid | 已付款 |
| abandon | 已放弃入职 |
| proposal_active | 应聘职位 |
| interview_accept | 接受面试 |
| interview_refuse | 拒绝面试 |
| review_accept | 接受二面 |
| review_refuse | 拒绝二面 |
| offer_accept | 接受offer |
| offer_refuse | 拒绝offer |
