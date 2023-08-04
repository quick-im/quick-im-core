package config

type msgGroupTopic string

const (
	MqMsgPrefix msgGroupTopic = "quickim.msg.*"
	// 消息持久化组件加入同一个订阅组，随机一个进行消费
	MqMsgPersistenceGroup msgGroupTopic = "quickim.msg.persistence"
	// 消息网关单独订阅一个主题，每一个网关都接受消息
	MqMsgConversationSubject msgGroupTopic = "quickim.msg.conversation"
)
