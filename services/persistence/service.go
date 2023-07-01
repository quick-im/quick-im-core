package persistence

import (
	"fmt"

	"github.com/quick-im/quick-im-core/internal/tracing"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/smallnest/rpcx/server"
	"go.opentelemetry.io/otel"
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
	tracer, ctx, err := tracing.InitJaeger(s.serviceName, s.agentHostPort)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Jaeger: %v", err))
	}
	defer tracer.Shutdown(ctx)
	currentTrace := otel.Tracer(s.serviceName)
	plugin := plugin.NewServerTracingPlugin(currentTrace)
	ser.Plugins.Add(plugin)
	// _ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_GENERATE_MESSAGE_ID, s.GenerateMessageID, "")
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
}
