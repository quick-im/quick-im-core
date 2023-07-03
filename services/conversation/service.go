package conversation

import (
	"context"
	"fmt"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
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
	// ser.Plugins.Add(opentracingrpc.ServerPlugin(tracer))
	ctx := context.Background()
	ctx = context.WithValue(ctx, contant.CTX_POSTGRES_KEY, db.GetDb())
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_CREATE_CONVERSATION, s.CreateConvercation(ctx), "")
	_ = ser.RegisterFunctionName(SERVER_NAME, SERVICE_JOIN_CONVERSATION, s.JoinConvercation(ctx), "")
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
}
