package persistence

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
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
	nc := s.InitNats()
	defer nc.Close()
	go s.listenMsg(ctx, nc)
	// 在服务端添加 Jaeger 拦截器
	if s.GetOpenTracing() {
		tracer, ctx := plugin.AddServerTrace(ser, SERVER_NAME, s.GetJeagerAgentHostPort())
		defer tracer.Shutdown(ctx)
	}
	s.addRegistryPlugin(ser)
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_SAVE_MSG_TO_DB, s.SaveMsgToDb(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_LAST30_MSG_FROM_DB, s.GetLast30MsgFromDb(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_MSG_FROM_DB_IN_RANGE, s.GetMsgFromDbInRange(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_THE_30MSG_AFTER_THE_ID, s.GetThe30MsgAfterTheId(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_THE_30MSG_BEFORE_THE_ID, s.GetThe30MsgBeforeTheId(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_LASTONE_MSG, s.GetLastOneMsgFromDb(ctx), "")
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.GetIp(), s.GetPort()))
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
