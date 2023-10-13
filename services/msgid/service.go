package msgid

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
	if s.GetOpenTracing() {
		tracer, ctx := plugin.AddServerTrace(ser, SERVER_NAME, s.GetJeagerAgentHostPort())
		defer tracer.Shutdown(ctx)
	}
	s.addRegistryPlugin(ser)
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GENERATE_MESSAGE_ID, s.GenerateMessageID(ctx), "")
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
