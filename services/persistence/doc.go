/*
	消息持久化服务模块
	Dep：
		Db
		Nats
	Feature：
		- 消息持久化存储
		- 离线消息收件箱机制处理
			- 更新ConversationID的lastMsgId
*/

package persistence
