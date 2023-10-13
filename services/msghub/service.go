package msghub

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/quick-im/quick-im-core/services/msgbroker"
	"github.com/quick-im/quick-im-core/services/persistence"
	cserver "github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/server"
)

type rpcxServer struct {
	config.ServiceConfig
}

func NewServer(opts ...config.Option) *rpcxServer {
	return &rpcxServer{
		config.NewServer(SERVER_NAME, opts...),
	}
}

func (s *rpcxServer) Start(ctx context.Context) error {
	ser := server.NewServer()
	// 在服务端添加 Jaeger 拦截器
	if s.GetOpenTracing() {
		tracer, ctx := plugin.AddServerTrace(ser, SERVER_NAME, s.GetJeagerAgentHostPort())
		defer tracer.Shutdown(ctx)
	}
	s.addRegistryPlugin(ser)
	nc := s.InitNats()
	defer nc.Close()
	ctx = context.WithValue(ctx, contant.CTX_NATS_KEY, nc)
	persistence := s.InitDepServices(persistence.SERVER_NAME)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_PERSISTENCE, persistence)
	defer persistence.CloseAndShutdownTrace()
	msgbroker := s.InitDepServices(msgbroker.SERVER_NAME)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_MSGBORKER, msgbroker)
	defer msgbroker.CloseAndShutdownTrace()
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_SEND_MSG, s.SendMsg(ctx), "")
	// s.logger.Info(s.serviceName, fmt.Sprintf("start at %s:%d", s.ip, s.port))
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.GetIp(), s.GetPort()))
}

func (s *rpcxServer) addRegistryPlugin(ser *server.Server) {
	if !s.GetUseConsulRegistry() {
		return
	}
	r := &cserver.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + fmt.Sprintf("%s:%d", s.GetIp(), s.GetPort()),
		ConsulServers:  s.GetConsulServers(),
		BasePath:       config.ServerPrefix,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	ser.Plugins.Add(r)
}

func (s *rpcxServer) InitNats() *messaging.NatsWarp {
	nc := messaging.NewNatsWithOpt(
		messaging.WithServers(s.GetNatsServers()...),
		messaging.WithJetStream(s.GetNatsEnableJetstream()),
	).GetNats()
	if s.GetNatsEnableJetstream() {
		js, err := nc.JetStream()
		if err != nil {
			s.GetLogger().Fatal("get nats jetstream err", fmt.Sprintf("%v", err))
		}
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     config.NatsStreamName,
			Subjects: []string{config.MqMsgPrefix},
		})
		if err != nil {
			s.GetLogger().Fatal("add stream to nats jetstream err", fmt.Sprintf("%v", err))
		}
	}
	return nc
}

func (r *rpcxServer) InitDepServices(serviceName string) *rpcx.RpcxClientWithOpt {
	service, err := rpcx.NewClient(
		rpcx.WithBasePath(config.ServerPrefix),
		rpcx.WithUseConsulRegistry(r.GetUseConsulRegistry()),
		rpcx.WithConsulServers(r.GetConsulServers()...),
		rpcx.WithServiceName(serviceName),
		rpcx.WithClientName(SERVER_NAME),
		rpcx.WithOpenTracing(r.GetOpenTracing()),
		rpcx.WithJeagerAgentHostPort(r.GetJeagerAgentHostPort()),
	)
	if err != nil {
		r.GetLogger().Fatal("init dep err", fmt.Sprintf("serviceName: %s , Err: %v", serviceName, err))
	}
	return service
}
