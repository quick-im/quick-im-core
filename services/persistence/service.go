package persistence

import (
	"fmt"

	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/smallnest/rpcx/server"
)

type rpcxServer struct {
	ip            string
	port          uint16
	openTracing   bool
	serviceName   string
	agentHostPort string
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
		tracer, ctx := plugin.AddServerTrace(ser, s.serviceName, s.agentHostPort)
		defer tracer.Shutdown(ctx)
	}
	// _ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GENERATE_MESSAGE_ID, s.GenerateMessageID, "")
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
}
