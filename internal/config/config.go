package config

import (
	"net"
)

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

type Config struct {
	Logger struct {
		Level int8
		Path  string
	}
	Jaeger struct {
		Host string
		Port uint16
	}
	PGDb struct {
		Host     string
		Port     uint16
		Username string
		Password string
		DbName   string
	}
	RethinkDb struct {
		Servers  []string
		Db       string
		Authkey  string
		Username string
		Password string
	}
	RedisDb struct {
		Host     string
		Port     uint16
		Username string
		Password string
	}
	Nats struct {
		Servers []string
	}
	Service struct {
		Conversation struct {
			Ip                string
			Port              uint16
			OpenTracing       bool
			UseConsulRegistry bool
		}
		Msgbroker struct {
			Ip                string
			Port              uint16
			OpenTracing       bool
			UseConsulRegistry bool
		}
		MsgId struct {
			Ip                string
			Port              uint16
			OpenTracing       bool
			UseConsulRegistry bool
		}
		Persistence struct {
			Ip                string
			Port              uint16
			OpenTracing       bool
			UseConsulRegistry bool
		}
	}
	Application struct {
		Ip                string
		Port              uint16
		OpenTracing       bool
		UseConsulRegistry bool
		IpWrite           []string
	}
}
