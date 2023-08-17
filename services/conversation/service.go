package conversation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/logger/innerzap"
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
	// 在服务端添加 Jaeger 拦截器
	if s.openTracing {
		tracer, ctx := plugin.AddServerTrace(ser, s.serviceName, s.trackAgentHostPort)
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
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
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
