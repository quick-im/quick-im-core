package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/logger/innerzap"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	cserver "github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap/zapcore"
)

type rpcxServer struct {
	ip                 string
	port               uint16
	openTracing        bool
	serviceName        string
	trackAgentHostPort string
	useConsulRegistry  bool
	consulServers      []string
	natsServers        []string
	logger             logger.Logger
}

func NewServer(opts ...Option) *rpcxServer {
	ser := &rpcxServer{}
	for i := range opts {
		opts[i](ser)
	}
	if ser.logger == nil {
		ser.logger = innerzap.NewZapLoggerAdapter(
			innerzap.NewLoggerWithOpt(
				innerzap.WithLogLevel(zapcore.DebugLevel),
				innerzap.WithServiceName(SERVER_NAME),
				innerzap.WithLogPath("logs"),
			).NewLogger(),
		)
	}
	return ser
}

func (s *rpcxServer) Start() error {
	ser := server.NewServer()
	// 在服务端添加 Jaeger 拦截器
	if s.openTracing {
		tracer, ctx := plugin.AddServerTrace(ser, s.serviceName, s.trackAgentHostPort)
		defer tracer.Shutdown(ctx)
	}
	ctx := context.Background()
	dbOpt := db.NewPostgresWithOpt(
		db.WithHost("localhost"),
		db.WithPort(5432),
		db.WithUsername("postgres"),
		db.WithPassword("123456"),
		db.WithDbName("quickim"),
	)
	ctx = context.WithValue(ctx, contant.CTX_POSTGRES_KEY, dbOpt.GetDb())
	s.addRegistryPlugin(ser)
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_SAVE_MSG_TO_DB, s.SaveMsgToDb(ctx), "")
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
}

func (s *rpcxServer) addRegistryPlugin(ser *server.Server) {
	if !s.useConsulRegistry {
		return
	}
	r := &cserver.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + fmt.Sprintf("%s:%d", s.ip, s.port),
		ConsulServers:  s.consulServers,
		BasePath:       SERVER_NAME,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	ser.Plugins.Add(r)
}
