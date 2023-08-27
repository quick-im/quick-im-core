package msgbroker

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/logger/innerzap"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/quick-im/quick-im-core/services/conversation"
	cserver "github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/server"
	"go.uber.org/zap/zapcore"
)

// type connList struct {
// 	lock    sync.RWMutex
// 	connMap map[string]connInfo
// }

// type connInfo struct {
// 	PlatformConn map[uint8]net.Conn
// 	Uid          string
// 	SessionId    string
// }

// TODO: 考虑将clientAddr更换为网关节点的uuid，防止冲突

type clientList struct {
	lock sync.RWMutex
	// 每个msgbroker可以接入若干gateway节点，同一个session不同platform可能接入不同gateway节点，所以这里做一下区分
	sessonIndex map[string]map[uint8]string // map[{sessionId}][{platform}]{clientAddr=>uuid}
	client      map[string]clientInfo       // map[{clientAddr=>uuid}]clientInfo
}

type clientInfo struct {
	conn    net.Conn
	connMap map[string]map[uint8]struct{}
}

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
	// connList            connList
	clientList clientList
}

func NewServer(opts ...Option) *rpcxServer {
	ser := &rpcxServer{
		consulServers:       make([]string, 0),
		natsServers:         make([]string, 0),
		natsEnableJetstream: true,
		serviceName:         SERVER_NAME,
		// connList: connList{
		// 	lock:    sync.RWMutex{},
		// 	connMap: make(map[string]connInfo, 100),
		// },
		clientList: clientList{
			lock:        sync.RWMutex{},
			sessonIndex: map[string]map[uint8]string{},
			client:      map[string]clientInfo{},
		},
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
	// 在服务端添加 Jaeger 拦截器
	if s.openTracing {
		tracer, ctx1 := plugin.AddServerTrace(ser, s.serviceName, s.trackAgentHostPort)
		defer tracer.Shutdown(ctx1)
	}
	conversationService := s.InitDepServices(conversation.SERVER_NAME)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_CONVERSATION, conversationService)
	defer conversationService.CloseAndShutdownTrace()
	selfService := s.InitDepServices(SERVER_NAME)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_MSGBORKER, selfService)
	defer selfService.CloseAndShutdownTrace()
	go s.listenMsg(ctx, nc)
	s.addRegistryPlugin(ser)
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_BROADCAST_RECV, s.BroadcastRecv(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_REGISTER_SESSION, s.RegisterSession(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_KICKOUT_DUPLICATE, s.KickoutDuplicate(ctx), "")
	// s.logger.Info(s.serviceName, fmt.Sprintf("start at %s:%d", s.ip, s.port))
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
}

func (s *rpcxServer) InitNats() *messaging.NatsWarp {
	nc := messaging.NewNatsWithOpt(
		messaging.WithServers(s.natsServers...),
		messaging.WithJetStream(s.natsEnableJetstream),
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

func (r *rpcxServer) InitDepServices(serviceName string) *rpcx.RpcxClientWithOpt {
	service, err := rpcx.NewClient(
		rpcx.WithBasePath(config.ServerPrefix),
		rpcx.WithUseConsulRegistry(r.useConsulRegistry),
		rpcx.WithConsulServers(r.consulServers...),
		rpcx.WithServiceName(serviceName),
		rpcx.WithClientName(r.serviceName),
		rpcx.WithOpenTracing(r.openTracing),
		rpcx.WithJeagerAgentHostPort(r.trackAgentHostPort),
	)
	if err != nil {
		r.logger.Fatal("init dep err", fmt.Sprintf("serviceName: %s , Err: %v", serviceName, err))
	}
	return service
}
