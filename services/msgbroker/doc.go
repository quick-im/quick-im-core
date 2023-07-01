/*
	消息分发/路由服务模块
	Dep：
		Nats
		GatewaySvc
	Feature：
		- 在线消息分发
			- 获取长连接服务器接入Session列表，将消息传递到指定服务器
		- 离线消息收件箱机制
			- 如果指定RecvSession于接收消息的过程中下线，则更新该Session对应的ConversationId的lastRecvMsgID
*/

package msgbroker
