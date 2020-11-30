namespace go uims_rpc_api

// 返回结构体
struct Response {
    1: string status; // 返回状态码 success | failed
    2: string msg;    // 状态码对应的提示语
    3: string data;   // 业务数据
}

// 服务体
service UIMSRpcApiService {
    // InvokeMethod UIMS系统为其它系统提供的接口通用调用方法
    // json string 参数
    Response InvokeMethod(1:string params)
}
