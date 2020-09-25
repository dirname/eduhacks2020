package routers

const (
	APIManagerClientGet   = "/api/manager/client/get"         // 管理端获取连接的客户端信息
	APIManagerRegister    = "/api/manager/register"           // 注册管理端系统的接口
	APIManagerMsgSend     = "/api/manager/client/send"        // 管理端发送消息到客户端
	APIManagerCloseClient = "/api/manager/client/close"       // 管理端关闭一个客户端连接的接口
	APIManagerMsgSends    = "/api/manager/clients/send"       // 批量发送
	APIManagerGroupSend   = "/api/manager/clients/group/send" // 发送群组
	APIManagerGroupBind   = "/api/manager/clients/group/bind" // 绑定群组
	APIManagerGroupGet    = "/api/manager/clients/group/get"  // 获取群组
)
