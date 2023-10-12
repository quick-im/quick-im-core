package conversation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/quick-im/quick-im-core/internal/config"
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
	// nc := s.InitNats()
	// defer nc.Close()
	// 是否使用消息队列待定
	// go s.listenMsg(ctx, nc)
	// 在服务端添加 Jaeger 拦截器
	if s.GetOpenTracing() {
		tracer, ctx := plugin.AddServerTrace(ser, SERVER_NAME, s.GetJeagerAgentHostPort())
		defer tracer.Shutdown(ctx)
	}
	// ser.Plugins.Add(opentracingrpc.ServerPlugin(tracer))
	s.addRegistryPlugin(ser)
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_CREATE_CONVERSATION, s.CreateConversation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_JOIN_CONVERSATION, s.JoinConversation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_ARCHIVE_CONVERSATIONS, s.SetArchiveConversations(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_DELETE_CONVERSATIONS, s.SetDeleteConversation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_CHECK_JOINED_CONVERSATION, s.CheckJoinedConversation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_CONVERSATION_INFO, s.GetConversationInfo(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_JOINED_CONVERSATIONS, s.GetJoinedConversations(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_KICKOUT_FOR_CONVERSATION, s.KickoutForConversation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_UPDATE_CONVERSATION_LASTMSG, s.UpdateConversationLastMsg(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_CONVERSATION_SSESSIONS, s.GetConversationSessions(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_UPDATE_SESSIONS_LAST_MSG, s.UpdateSessionLastRecvMsg(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_LASTONE_MSGID_FROM_DB, s.GetLastOneMsgIdFromDb(ctx), "")
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.GetIp(), s.GetPort()))
}

// func (s *rpcxServer) InitNats() *messaging.NatsWarp {
// 	nc := messaging.NewNatsWithOpt(
// 		messaging.WithServers(s.natsServers...),
// 	).GetNats()
// 	if s.natsEnableJetstream {
// 		js, err := nc.JetStream()
// 		if err != nil {
// 			s.logger.Fatal("get nats jetstream err", fmt.Sprintf("%v", err))
// 		}
// 		_, err = js.AddStream(&nats.StreamConfig{
// 			Name:     config.NatsStreamName,
// 			Subjects: []string{config.MqMsgPrefix},
// 		})
// 		if err != nil {
// 			s.logger.Fatal("add stream to nats jetstream err", fmt.Sprintf("%v", err))
// 		}
// 	}
// 	return nc
// }

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
