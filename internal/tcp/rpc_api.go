package tcp

import (
	"uims/internal/thriftcontroller/cass/groupcontroller"
	"uims/internal/thriftcontroller/cass/usercontroller"
	"uims/internal/thriftcontroller/democontroller"
	"uims/internal/thriftcontroller/esigncontroller"
	"uims/internal/thriftcontroller/oauthcontroller"
	usercontroller2 "uims/internal/thriftcontroller/vzhuo/usercontroller"
	"uims/internal/thriftcontroller/vzhuo/wechatcontroller"
	thriftserver "uims/pkg/thrift/server"
)

var ThriftRPCmethodMap = thriftserver.MethodMap{
	"test": democontroller.Demo,
	// 获取一个唯一ID (暂未通过业务系统端的单元测试，不能使用)
	//"uuid": uuidcontroller.Generate,
	"save_cass_push_user_data":  usercontroller.Create,
	"save_cass_push_group_data": groupcontroller.Create,
	// 保存微桌用户推送的个人资料
	"save_vzhuo_push_user_data": usercontroller2.UpdateUserInfo,
	// 保存微桌用户上传的身份认证数据
	"save_vzhuo_user_identity": usercontroller2.UpdateUserIdentity,
	// 通过 code 获取 access_token
	"sns/oauth/access_token": oauthcontroller.AccessToken,
	// 通过 refresh 更新 access_token
	"sns/oauth/refresh_token": oauthcontroller.RefreshToken,
	// 通过 access_token 获取用户信息
	"sns/userinfo": oauthcontroller.UserInfo,
	// 获取绑定微信 url
	"sns/bind/wechat/url":   oauthcontroller.GetBindWeChatURL,
	"sns/unbind/wechat/url": oauthcontroller.GetUnbindWeChatURL,
	"sns/isbind/wechat":     oauthcontroller.IsBindWeChat,
	// 微桌微信关注事件
	"vzhuo/wechat/follow": wechatcontroller.FollowEvent,
	// 微桌后台重置密码
	"vzhuo/reset/password": usercontroller2.UpdateUserPassword,
	// 保存微桌代发项目生成的user
	"vzhuo/save/user": usercontroller2.SaveVzhuoUser,
	// 电签完成通知
	"esign/notify": esigncontroller.NotifyESign,
}
