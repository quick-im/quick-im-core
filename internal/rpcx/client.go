package rpcx

import (
	"context"
	"errors"
	"fmt"

	"github.com/quick-im/quick-im-core/internal/quickim_errors"
	"github.com/quick-im/quick-im-core/internal/tracing"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/smallnest/rpcx/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type RpcxClientWithOpt struct {
	openTracing              bool
	useConsulRegistry        bool
	serviceName              string
	serverAddress            string
	consulServers            []string
	trackJaegeragentHostPort string
	tracePtr                 *trace.TracerProvider
	xclient                  client.XClient
	ctx                      context.Context
}

type ctxInitInner string

type rpcxOption func(*RpcxClientWithOpt)

func WithOpenTracing(disable bool) rpcxOption {
	return func(rs *RpcxClientWithOpt) {
		rs.openTracing = disable
	}
}

func WithServiceName(serviceName string) rpcxOption {
	return func(rs *RpcxClientWithOpt) {
		rs.serviceName = serviceName
	}
}

func WithJeagerAgentHostPort(JaegeragentHostPort string) rpcxOption {
	return func(rs *RpcxClientWithOpt) {
		rs.trackJaegeragentHostPort = JaegeragentHostPort
	}
}

func WithServerAddress(server string) rpcxOption {
	return func(rco *RpcxClientWithOpt) {
		rco.serverAddress = server
	}
}

func WithUseConsulRegistry(useConsulRegistry bool) rpcxOption {
	return func(rco *RpcxClientWithOpt) {
		rco.useConsulRegistry = useConsulRegistry
	}
}

func WithConsulServer(server string) rpcxOption {
	return func(rco *RpcxClientWithOpt) {
		rco.consulServers = append(rco.consulServers, server)
	}
}

func WithConsulServers(servers ...string) rpcxOption {
	return func(rco *RpcxClientWithOpt) {
		rco.consulServers = servers
	}
}

func NewClient(opts ...rpcxOption) (*RpcxClientWithOpt, error) {
	c := &RpcxClientWithOpt{
		consulServers: make([]string, 0),
	}
	for i := range opts {
		opts[i](c)
	}
	d, err := client.NewPeer2PeerDiscovery(c.serverAddress, "")
	if err != nil {
		return nil, err
	}
	xclient := client.NewXClient(c.serviceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	if c.openTracing {
		tracer, ctx := c.addTrace(xclient)
		c.tracePtr = tracer
		c.ctx = ctx
	}
	c.xclient = xclient
	return c, nil
}

func (s *RpcxClientWithOpt) Close() error {
	return s.xclient.Close()
}

func (s *RpcxClientWithOpt) ShutdownTrace() error {
	if s.tracePtr == nil {
		return quickim_errors.ErrTraceClosed
	}
	return s.tracePtr.Shutdown(s.ctx)
}

func (s *RpcxClientWithOpt) CloseAndShutdownTrace() error {
	var err error
	err1 := s.xclient.Close()
	if err1 != nil {
		err = err1
	}
	if s.tracePtr != nil {
		err2 := s.tracePtr.Shutdown(s.ctx)
		if err2 != nil {
			err = errors.Join(err, err2)
		}
	}
	return err
}

func (s *RpcxClientWithOpt) Call(ctx context.Context, serviceMethod string, arg interface{}, replay interface{}) error {
	ctxInner := context.WithValue(s.ctx, ctxInitInner("initCtx"), nil)
	return s.xclient.Call(ctxInner, serviceMethod, arg, replay)
}

func (s *RpcxClientWithOpt) addTrace(xclient client.XClient) (*trace.TracerProvider, context.Context) {
	// 添加 Jaeger 拦截器
	plugins := client.NewPluginContainer()
	tracer, ctx, err := tracing.InitJaeger("client", s.trackJaegeragentHostPort)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Jaeger: %v", err))
	}
	ts := otel.Tracer("client")
	plugins.Add(plugin.NewClientTracingPlugin(ts))
	xclient.SetPlugins(plugins)
	return tracer, ctx
}
