package config

import "net"

const (
	// 接口通信密钥
	ServiceKey = "quick-im"
	// Nats流名称
	NatsStreamName = "QUICKIM_MSG_STREAM"
	// 消息订阅主题前缀
	MqMsgPrefix = "quickim.msg.>"
	// 消息持久化组件加入同一个订阅组，随机一个进行消费
	MqMsgPersistenceGroup = "quickim.msg.persistence"
	// 消息网关单独订阅一个主题，每一个网关都接受消息
	MqMsgBrokerSubject = "quickim.msg.msgbroker"
	// 消息持久化表
	RethinkMsgDb = "msg"
	// 注册中心服务前缀
	ServerPrefix = "quick.im.instance.1"
	// TLS证书
	PublicCert  = "cert/server.crt"
	PriviteCert = "cert/server.key"
)

var (
	// 调试使用，生产环境不建议
	_, all, _ = net.ParseCIDR("0.0.0.0/0")
	IPWhite   = []*net.IPNet{
		all,
	}
)
