package conversation

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	cserver "github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/server"
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
}

func NewServer(opts ...Option) *rpcxServer {
	ser := &rpcxServer{}
	for i := range opts {
		opts[i](ser)
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
	// ser.Plugins.Add(opentracingrpc.ServerPlugin(tracer))
	ctx := context.Background()
	dbOpt := db.NewPostgresWithOpt(
		db.WithHost("localhost"),
		db.WithPort(5432),
		db.WithUsername("postgres"),
		db.WithPassword("123456"),
		db.WithDbName("quickim"),
	)
	ctx = context.WithValue(ctx, contant.CTX_POSTGRES_KEY, dbOpt.GetDb())
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_CREATE_CONVERSATION, s.CreateConvercation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_JOIN_CONVERSATION, s.JoinConvercation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_ARCHIVE_CONVERCATIONS, s.SetArchiveConversations(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_DELETE_CONVERCATIONS, s.SetDeleteConversation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_CHECK_JOINED_CONVERCATION, s.CheckJoinedConversation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_CONVERCATION_INFO, s.GetConversationInfo(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GET_JOINED_CONVERCATIONS, s.GetJoinedConversations(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_KICKOUT_FOR_CONVERCATION, s.KickoutForConversation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_UPDATE_CONVERCATIONS_LASTMSG, s.UpdateConversationLastMsg(ctx), "")
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
