package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/logger/innerzap"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
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
		serviceName:         SERVER_NAME,
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

func (s *rpcxServer) Start(ctx context.Context) error {
	ser := server.NewServer()
	nc := s.InitNats()
	defer nc.Close()
	go s.listenMsg(nc)
	// 在服务端添加 Jaeger 拦截器
	if s.openTracing {
		tracer, ctx := plugin.AddServerTrace(ser, s.serviceName, s.trackAgentHostPort)
		defer tracer.Shutdown(ctx)
	}
	s.addRegistryPlugin(ser)
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_SAVE_MSG_TO_DB, s.SaveMsgToDb(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_LAST30_MSG_FROM_DB, s.GetLast30MsgFromDb(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_MSG_FROM_DB_IN_RANGE, s.GetMsgFromDbInRange(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_THE_30MSG_AFTER_THE_ID, s.GetThe30MsgAfterTheId(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_THE_30MSG_BEFORE_THE_ID, s.GetThe30MsgBeforeTheId(ctx), "")
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
}

func (s *rpcxServer) InitNats() *messaging.NatsWarp {
	nc := messaging.NewNatsWithOpt(
		messaging.WithServers(s.natsServers...),
	).GetNats()
	if s.natsEnableJetstream {
		js, err := nc.JetStream()
		if err != nil {
			s.logger.Fatal("get nats jetstream err", fmt.Sprintf("%v", err))
		}
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     config.NatsStreamName,
			Subjects: []string{config.MqMsgPrefix},
		})
		if err != nil {
			s.logger.Fatal("add stream to nats jetstream err", fmt.Sprintf("%v", err))
		}
	}
	return nc
}

func (s *rpcxServer) addRegistryPlugin(ser *server.Server) {
	if !s.useConsulRegistry {
		return
	}
	r := &cserver.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + fmt.Sprintf("%s:%d", s.ip, s.port),
		ConsulServers:  s.consulServers,
		BasePath:       config.ServerPrefix,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	ser.Plugins.Add(r)
}
