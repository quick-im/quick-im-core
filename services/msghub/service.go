package msghub

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/logger/innerzap"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/quick-im/quick-im-core/services/persistence"
	cserver "github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap/zapcore"
)

type rpcxServer struct {
	ip                  string
	port                uint16
	openTracing         bool
	serviceName         string
	trackAgentHostPort  string
	useConsulRegistry   bool
	consulServers       []string
	natsServers         []string
	natsEnableJetstream bool
	logger              logger.Logger
}

func NewServer(opts ...Option) *rpcxServer {
	ser := &rpcxServer{
		consulServers:       make([]string, 0),
		natsServers:         make([]string, 0),
		natsEnableJetstream: true,
	}
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
	s.addRegistryPlugin(ser)
	ctx := context.Background()
	nc := s.InitNats()
	defer nc.Close()
	ctx = context.WithValue(ctx, contant.CTX_NATS_KEY, nc)
	persistence := s.InitDepServices(persistence.SERVER_NAME)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_PERSISTENCE, persistence)
	defer persistence.CloseAndShutdownTrace()
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_SEND_MSG, s.SendMsg(ctx), "")
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

func (s *rpcxServer) InitNats() *nats.Conn {
	nc := messaging.NewNatsWithOpt(
		messaging.WithServers(s.natsServers...),
	).GetNats()
	if s.natsEnableJetstream {
		js, err := nc.JetStream()
		if err != nil {
			s.logger.Fatal("get nats jetstream err", fmt.Sprintf("%v", err))
		}
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     "MSG_STREAM",
			Subjects: []string{"stream.msg.>"},
		})
		if err != nil {
			s.logger.Fatal("add stream to nats jetstream err", fmt.Sprintf("%v", err))
		}
	}
	return nc
}

func (r *rpcxServer) InitDepServices(serviceName string) *rpcx.RpcxClientWithOpt {
	service, err := rpcx.NewClient(
		rpcx.WithUseConsulRegistry(r.useConsulRegistry),
		rpcx.WithConsulServers(r.consulServers...),
		rpcx.WithServiceName(serviceName),
		rpcx.WithOpenTracing(r.openTracing),
		rpcx.WithJeagerAgentHostPort(r.trackAgentHostPort),
	)
	if err != nil {
		r.logger.Fatal("init dep %s err", serviceName, fmt.Sprintf("%v", err))
	}
	return service
}
