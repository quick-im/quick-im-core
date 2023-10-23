package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
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
	// 服务名称
	ConversationSerName = "conversation"
	MsgbrokerSerName    = "msgbroker"
	MsghubSerName       = "msghub"
	MsgIdSerName        = "msgid"
	PersistenceSerName  = "persistence"
	GatewayName         = "gateway"
)

var (
// 调试使用，生产环境不建议
// _, all, _ = net.ParseCIDR("0.0.0.0/0")
//
//	IPWhite   = []*net.IPNet{
//		all,
//	}
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

var GatewayIpwriteFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    "gateway.ipWrite",
		Value:   cli.NewStringSlice("0.0.0.0/0"),
		Usage:   "gateway ipwrite",
		EnvVars: []string{"IP_WRITE"},
	},
	&cli.StringFlag{
		Name:     "gateway.jwtKey",
		Value:    "quickimkey",
		Usage:    "gateway jwt key",
		Required: true,
		EnvVars:  []string{"GATEWAY_JWT_KEY"},
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
		Value:   "123456",
		Usage:   "postgres password",
		EnvVars: []string{"PG_PASSWORD"},
	},
	&cli.StringFlag{
		Name:    "pg.dbname",
		Value:   "quickim",
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

func GetServiceFlags(service string, port uint) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    fmt.Sprintf("%s.ip", service),
			Value:   "0.0.0.0",
			Usage:   "listen ip",
			EnvVars: []string{"IP"},
		},
		&cli.UintFlag{
			Name:    fmt.Sprintf("%s.port", service),
			Value:   port,
			Usage:   "listen port",
			EnvVars: []string{"PORT"},
		},
	}
}

type IMConf struct {
	OpenTracing bool      `toml:"openTracing"`
	Logger      Logger    `toml:"logger"`
	Jaeger      Jaeger    `toml:"jaeger"`
	Postgres    Postgres  `toml:"postgres"`
	Redis       Redis     `toml:"redis"`
	Nats        Nats      `toml:"nats"`
	RethinkDb   RethinkDb `toml:"rethinkdb"`
	Consul      Consul    `toml:"consul"`
	Gateway     Gateway   `toml:"gateway"`
	Services    Services  `toml:"services"`
}
type Logger struct {
	Path  string `toml:"path"`
	Level int    `toml:"level"`
}
type Jaeger struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}
type Postgres struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Db       string `toml:"db"`
}
type Redis struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}
type Nats struct {
	Servers []string `toml:"servers"`
}
type RethinkDb struct {
	Servers  []string `toml:"servers"`
	Db       string   `toml:"db"`
	Authkey  string   `toml:"authkey"`
	Username string   `toml:"username"`
	Password string   `toml:"password"`
}
type Consul struct {
	Servers []string `toml:"servers"`
}
type Gateway struct {
	IP      string   `toml:"ip"`
	Port    int      `toml:"port"`
	IPWrite []string `toml:"ipWrite"`
	JwtKey  string   `toml:"jwtKey"`
}
type Conversation struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}
type Msgbroker struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}
type Msghub struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}
type Msgid struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}
type Persistence struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}
type Services struct {
	Conversation Conversation `toml:"conversation"`
	Msgbroker    Msgbroker    `toml:"msgbroker"`
	Msghub       Msghub       `toml:"msghub"`
	Msgid        Msgid        `toml:"msgid"`
	Persistence  Persistence  `toml:"persistence"`
}

func MergeConf(args *cli.Context) IMConf {
	var conf IMConf
	readFile := false
	if args.IsSet("config") {
		_, err := toml.DecodeFile(args.String("config"), &conf)
		if err != nil {
			println(err.Error())
		} else {
			readFile = true
		}
	}
	conf.OpenTracing = merge(args.IsSet("jaeger.enable") || !readFile, conf.OpenTracing, args.Bool("jaeger.enable"))
	conf.Logger.Path = merge(args.IsSet("log.path") || !readFile, conf.Logger.Path, args.String("log.path"))
	conf.Logger.Level = merge(args.IsSet("log.level") || !readFile, conf.Logger.Level, args.Int("log.level"))
	conf.Jaeger.Host = merge(args.IsSet("jaeger.host") || !readFile, conf.Jaeger.Host, args.String("jaeger.host"))
	conf.Jaeger.Port = merge(args.IsSet("jaeger.port") || !readFile, conf.Jaeger.Port, args.Int("jaeger.port"))
	conf.Postgres.Host = merge(args.IsSet("pg.host") || !readFile, conf.Postgres.Host, args.String("pg.host"))
	conf.Postgres.Port = merge(args.IsSet("pg.port") || !readFile, conf.Postgres.Port, args.Int("pg.port"))
	conf.Postgres.Username = merge(args.IsSet("pg.username") || !readFile, conf.Postgres.Username, args.String("pg.username"))
	conf.Postgres.Password = merge(args.IsSet("pg.password") || !readFile, conf.Postgres.Password, args.String("pg.password"))
	conf.Postgres.Db = merge(args.IsSet("pg.dbname") || !readFile, conf.Postgres.Db, args.String("pg.dbname"))
	conf.Redis.Host = merge(args.IsSet("redis.host") || !readFile, conf.Redis.Host, args.String("redis.host"))
	conf.Redis.Port = merge(args.IsSet("redis.port") || !readFile, conf.Redis.Port, args.Int("redis.port"))
	conf.Redis.Username = merge(args.IsSet("redis.username") || !readFile, conf.Redis.Username, args.String("redis.username"))
	conf.Redis.Password = merge(args.IsSet("redis.password") || !readFile, conf.Redis.Password, args.String("redis.password"))
	conf.Nats.Servers = merge(args.IsSet("nats.servers") || !readFile, conf.Nats.Servers, args.StringSlice("nats.servers"))
	conf.RethinkDb.Servers = merge(args.IsSet("rethinkdb.servers") || !readFile, conf.RethinkDb.Servers, args.StringSlice("rethinkdb.servers"))
	conf.RethinkDb.Db = merge(args.IsSet("rethinkdb.db") || !readFile, conf.RethinkDb.Db, args.String("rethinkdb.db"))
	conf.RethinkDb.Authkey = merge(args.IsSet("rethinkdb.authkey") || !readFile, conf.RethinkDb.Authkey, args.String("rethinkdb.authkey"))
	conf.RethinkDb.Username = merge(args.IsSet("rethinkdb.username") || !readFile, conf.RethinkDb.Username, args.String("rethinkdb.username"))
	conf.RethinkDb.Password = merge(args.IsSet("rethinkdb.password") || !readFile, conf.RethinkDb.Password, args.String("rethinkdb.password"))
	conf.Consul.Servers = merge(args.IsSet("consul.servers") || !readFile, conf.Consul.Servers, args.StringSlice("consul.servers"))
	conf.Services.Conversation.IP = merge(args.IsSet("conversation.ip") || !readFile, conf.Services.Conversation.IP, args.String("conversation.ip"))
	conf.Services.Conversation.Port = merge(args.IsSet("conversation.port") || !readFile, conf.Services.Conversation.Port, args.Int("conversation.port"))
	conf.Services.Msgbroker.IP = merge(args.IsSet("msgbroker.ip") || !readFile, conf.Services.Msgbroker.IP, args.String("msgbroker.ip"))
	conf.Services.Msgbroker.Port = merge(args.IsSet("msgbroker.port") || !readFile, conf.Services.Msgbroker.Port, args.Int("msgbroker.port"))
	conf.Services.Msghub.IP = merge(args.IsSet("msghub.ip") || !readFile, conf.Services.Msghub.IP, args.String("msghub.ip"))
	conf.Services.Msghub.Port = merge(args.IsSet("msghub.port") || !readFile, conf.Services.Msghub.Port, args.Int("msghub.port"))
	conf.Services.Msgid.IP = merge(args.IsSet("msgid.ip") || !readFile, conf.Services.Msgid.IP, args.String("msgid.ip"))
	conf.Services.Msgid.Port = merge(args.IsSet("msgid.port") || !readFile, conf.Services.Msgid.Port, args.Int("msgid.port"))
	conf.Services.Persistence.IP = merge(args.IsSet("persistence.ip") || !readFile, conf.Services.Persistence.IP, args.String("persistence.ip"))
	conf.Services.Persistence.Port = merge(args.IsSet("persistence.port") || !readFile, conf.Services.Persistence.Port, args.Int("persistence.port"))
	conf.Gateway.IP = merge(args.IsSet("gateway.ip") || !readFile, conf.Gateway.IP, args.String("gateway.ip"))
	conf.Gateway.Port = merge(args.IsSet("gateway.port") || !readFile, conf.Gateway.Port, args.Int("gateway.port"))
	conf.Gateway.IPWrite = merge(args.IsSet("gateway.ipWrite") || !readFile, conf.Gateway.IPWrite, args.StringSlice("gateway.ipWrite"))
	conf.Gateway.JwtKey = merge(args.IsSet("gateway.jwtKey") || !readFile, conf.Gateway.JwtKey, args.String("gateway.jwtKey"))
	return conf
}

func merge[T any](set bool, src T, dest T) T {
	if set {
		return dest
	}
	return src
}

// type Config struct {
// 	Logger struct {
// 		Level int8
// 		Path  string
// 	}
// 	Jaeger struct {
// 		Enable bool // 打开这个则所有服务开启OpenTracing
// 		Host   string
// 		Port   uint16
// 	}
// 	PGDb struct {
// 		Host     string
// 		Port     uint16
// 		Username string
// 		Password string
// 		DbName   string
// 	}
// 	RethinkDb struct {
// 		Servers  []string
// 		Db       string
// 		Authkey  string
// 		Username string
// 		Password string
// 	}
// 	RedisDb struct {
// 		Host     string
// 		Port     uint16
// 		Username string
// 		Password string
// 	}
// 	Nats struct {
// 		Servers []string
// 	}
// 	Consul struct {
// 		Servers []string
// 	}
// 	Service struct {
// 		Conversation struct {
// 			Ip                string
// 			Port              uint16
// 			UseConsulRegistry bool
// 		}
// 		Msgbroker struct {
// 			Ip                string
// 			Port              uint16
// 			UseConsulRegistry bool
// 		}
// 		MsgId struct {
// 			Ip                string
// 			Port              uint16
// 			UseConsulRegistry bool
// 		}
// 		Persistence struct {
// 			Ip                string
// 			Port              uint16
// 			UseConsulRegistry bool
// 		}
// 	}
// 	Application struct {
// 		Ip                string
// 		Port              uint16
// 		UseConsulRegistry bool
// 		IpWrite           []string
// 	}
// }
