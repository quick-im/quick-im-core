package msgbroker

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/quick-im/quick-im-core/services/conversation"
	cserver "github.com/rpcxio/rpcx-consul/serverplugin"
	"github.com/smallnest/rpcx/server"
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
	config.ServiceConfig
	rpcxSer *server.Server
	// connList            connList
	clientList clientList
}

func NewServer(opts ...config.Option) *rpcxServer {
	ser := &rpcxServer{
		ServiceConfig: config.NewServer(SERVER_NAME, opts...),
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
	return ser
}

func (s *rpcxServer) Start(ctx context.Context) error {
	ser := server.NewServer()
	s.rpcxSer = ser
	nc := s.InitNats()
	defer nc.Close()
	// 在服务端添加 Jaeger 拦截器
	if s.GetOpenTracing() {
		tracer, ctx1 := plugin.AddServerTrace(ser, SERVER_NAME, s.GetJeagerAgentHostPort())
		defer tracer.Shutdown(ctx1)
	}
	conversationService := s.InitDepServices(conversation.SERVER_NAME)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_CONVERSATION, conversationService)
	defer conversationService.CloseAndShutdownTrace()
	selfService := s.InitDepServices(SERVER_NAME)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_MSGBORKER, selfService)
	defer selfService.CloseAndShutdownTrace()
	go s.listenMsg(ctx, nc)
	go s.Heartbeat(time.Minute)
	s.addRegistryPlugin(ser)
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_BROADCAST_RECV, s.BroadcastRecv(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_REGISTER_SESSION, s.RegisterSession(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_KICKOUT_DUPLICATE, s.KickoutDuplicate(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_UNREGISTER_SESSION, s.UnRegisterSession(ctx), "")
	// s.logger.Info(s.serviceName, fmt.Sprintf("start at %s:%d", s.ip, s.port))
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

func (s *rpcxServer) Heartbeat(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	c := codec.GobUtils[BroadcastMsgWarp]{}
	heartbeatData, err := c.Encode(BroadcastMsgWarp{
		Action: Heartbeat,
	})
	if err != nil {
		panic(err)
	}
	for range ticker.C {
		gateways := map[string]clientInfo{}
		s.clientList.lock.RLock()
		for gatewayUuid := range s.clientList.client {
			gateways[gatewayUuid] = s.clientList.client[gatewayUuid]
		}
		s.clientList.lock.RUnlock()
		// map[gateway]map[sessioId]map[platform]
		needGC := map[string]map[string]map[uint8]struct{}{}
		for gatewayUuid := range gateways {
			if err := s.rpcxSer.SendMessage(gateways[gatewayUuid].conn, SERVER_NAME, SERVICE_BROADCAST_RECV, nil, heartbeatData); err != nil {
				// s.logger.Error("Heartbeat Err:", fmt.Sprintf("gatewayUuid: %s, gatewayAddr: %s, err: %v", gatewayUuid, gateways[gatewayUuid].conn.RemoteAddr().String(), err))
				needGC[gatewayUuid] = gateways[gatewayUuid].connMap
			}
		}
		if len(needGC) > 0 {
			s.clientList.lock.Lock()
			for gatewayId := range needGC {
				for sessionId := range needGC[gatewayId] {
					for platform := range needGC[gatewayId][sessionId] {
						// 先检查key是否存在，防止其他地方清理之后导致的panic
						if _, sessionOk := s.clientList.sessonIndex[sessionId]; sessionOk {
							if g, platformOk := s.clientList.sessonIndex[sessionId][platform]; platformOk && g == gatewayId {
								if len(s.clientList.sessonIndex[sessionId]) == 1 {
									delete(s.clientList.sessonIndex, sessionId)
								} else {
									delete(s.clientList.sessonIndex[sessionId], platform)
								}
							}
							// 如果当前session在线平台只有一个且和需要清理的平台匹配，则直接删除这个session索引
						}
					}
				}
				delete(s.clientList.client, gatewayId)
			}
			s.clientList.lock.Unlock()
		}
	}
}
