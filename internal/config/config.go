package config

import (
	"net"

	"github.com/urfave/cli/v2"
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

var JaegerFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "jaeger.enable",
		Value:   false,
		Usage:   "jaeger enable",
		EnvVars: []string{"JAEGER_ENABLE"},
	},
	&cli.StringFlag{
		Name:    "jaeger.host",
		Value:   "127.0.0.1",
		Usage:   "jaeger host",
		EnvVars: []string{"JAEGER_HOST"},
	},
	&cli.UintFlag{
		Name:    "jaeger.port",
		Value:   6832,
		Usage:   "jaeger port",
		EnvVars: []string{"JAEGER_PORT"},
	},
}

var PgFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "pg.host",
		Value:   "127.0.0.1",
		Usage:   "postgres host",
		EnvVars: []string{"PG_HOST"},
	},
	&cli.UintFlag{
		Name:    "pg.port",
		Value:   5432,
		Usage:   "postgres port",
		EnvVars: []string{"PG_PORT"},
	},
	&cli.StringFlag{
		Name:    "pg.username",
		Value:   "postgres",
		Usage:   "postgres username",
		EnvVars: []string{"PG_USERNAME"},
	},
	&cli.StringFlag{
		Name:    "pg.password",
		Value:   "postgres",
		Usage:   "postgres password",
		EnvVars: []string{"PG_PASSWORD"},
	},
	&cli.StringFlag{
		Name:    "pg.dbname",
		Value:   "postgres",
		Usage:   "postgres dbname",
		EnvVars: []string{"PG_DBNAME"},
	},
}

var RedisFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "redis.host",
		Value:   "127.0.0.1",
		Usage:   "redis host",
		EnvVars: []string{"REDIS_HOST"},
	},
	&cli.UintFlag{
		Name:    "redis.port",
		Value:   6379,
		Usage:   "redis port",
		EnvVars: []string{"REDIS_PORT"},
	},
	&cli.StringFlag{
		Name:    "redis.username",
		Value:   "",
		Usage:   "redis username",
		EnvVars: []string{"REDIS_USERNAME"},
	},
	&cli.StringFlag{
		Name:    "redis.password",
		Value:   "",
		Usage:   "redis password",
		EnvVars: []string{"REDIS_PASSWORD"},
	},
}

var RethinkDbFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name: "rethinkdb.servers",
		Value: cli.NewStringSlice(
			"127.0.0.1:28015",
		),
		Usage:   "rethinkdb servers",
		EnvVars: []string{"RETHINKDB_SERVERS"},
	},
	&cli.StringFlag{
		Name:    "rethinkdb.db",
		Value:   "quickim",
		Usage:   "rethinkdb db name",
		EnvVars: []string{"RETHINKDB_DB"},
	},
	&cli.StringFlag{
		Name:    "rethinkdb.authkey",
		Value:   "",
		Usage:   "rethinkdb authkey",
		EnvVars: []string{"RETHINKDB_AUTHKEY"},
	},
	&cli.StringFlag{
		Name:    "rethinkdb.username",
		Value:   "",
		Usage:   "rethinkdb username",
		EnvVars: []string{"RETHINKDB_USERNAME"},
	},
	&cli.StringFlag{
		Name:    "rethinkdb.password",
		Value:   "",
		Usage:   "rethinkdb password",
		EnvVars: []string{"RETHINKDB_PASSWORD"},
	},
}

var NatsFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "nats.servers",
		Value:   cli.NewStringSlice("127.0.0.1:4222"),
		Usage:   "nats servers",
		EnvVars: []string{"NATS_SERVERS"},
	},
}

var ConsulFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "consul.servers",
		Value:   cli.NewStringSlice("127.0.0.1:8500"),
		Usage:   "consul servers",
		EnvVars: []string{"CONSUL_SERVERS"},
	},
}

var LogFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "log.path",
		Value:   "./logs/",
		Usage:   "日志文件路径",
		EnvVars: []string{"LOG_PATH"},
	},
	&cli.IntFlag{
		Name:  "log.level",
		Value: -1,
		Usage: `日志记录等级：
			-1：DebugLevel
			 1：InfoLevel
			 2：WarnLevel
			 3：ErrorLevel
			 4：DPanicLevel
			 5：PanicLevel
			 6：FatalLevel
			`,
		EnvVars: []string{"LOG_LEVEL"},
	},
}

type Config struct {
	Logger struct {
		Level int8
		Path  string
	}
	Jaeger struct {
		Enable bool // 打开这个则所有服务开启OpenTracing
		Host   string
		Port   uint16
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
	Consul struct {
		Servers []string
	}
	Service struct {
		Conversation struct {
			Ip                string
			Port              uint16
			UseConsulRegistry bool
		}
		Msgbroker struct {
			Ip                string
			Port              uint16
			UseConsulRegistry bool
		}
		MsgId struct {
			Ip                string
			Port              uint16
			UseConsulRegistry bool
		}
		Persistence struct {
			Ip                string
			Port              uint16
			UseConsulRegistry bool
		}
	}
	Application struct {
		Ip                string
		Port              uint16
		UseConsulRegistry bool
		IpWrite           []string
	}
}
