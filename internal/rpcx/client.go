package rpcx

import (
	"context"
	"fmt"

	"github.com/quick-im/quick-im-core/internal/errors"
	"github.com/quick-im/quick-im-core/internal/tracing"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/smallnest/rpcx/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type rpcxClientOpt struct {
	openTracing              bool
	useConsulRegistry        bool
	serviceName              string
	serverAddress            string
	consulServers            []string
	trackJaegeragentHostPort string
	tracePtr                 *trace.TracerProvider
	xclient                  client.XClient
}

type rpcxOption func(*rpcxClientOpt)

func WithOpenTracing(disable bool) rpcxOption {
	return func(rs *rpcxClientOpt) {
		rs.openTracing = disable
	}
}

func WithJeagerServiceName(serviceName string) rpcxOption {
	return func(rs *rpcxClientOpt) {
		rs.serviceName = serviceName
	}
}

func WithJeagerAgentHostPort(JaegeragentHostPort string) rpcxOption {
	return func(rs *rpcxClientOpt) {
		rs.trackJaegeragentHostPort = JaegeragentHostPort
	}
}

func WithServerAddress(server string) rpcxOption {
	return func(rco *rpcxClientOpt) {
		rco.serverAddress = server
	}
}

func WithUseConsulRegistry(useConsulRegistry bool) rpcxOption {
	return func(rco *rpcxClientOpt) {
		rco.useConsulRegistry = useConsulRegistry
	}
}

func WithConsulServer(server string) rpcxOption {
	return func(rco *rpcxClientOpt) {
		rco.consulServers = append(rco.consulServers, server)
	}
}

func WithConsulServers(servers ...string) rpcxOption {
	return func(rco *rpcxClientOpt) {
		rco.consulServers = servers
	}
}

func NewClient(opts ...rpcxOption) (*rpcxClientOpt, error) {
	c := &rpcxClientOpt{
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
		_ = ctx
	}
	c.xclient = xclient
	return c, nil
}

func (s *rpcxClientOpt) Close() error {
	return s.xclient.Close()
}

func (s *rpcxClientOpt) ShutdownTrace() error {
	if s.tracePtr == nil {
		return errors.ErrTraceClosed
	}
	return s.tracePtr.Shutdown(context.Background())
}

func (s *rpcxClientOpt) Call() error {
	return nil
}

func (s *rpcxClientOpt) addTrace(xclient client.XClient) (*trace.TracerProvider, context.Context) {
	// 添加 Jaeger 拦截器
	plugins := client.NewPluginContainer()
	tracer, ctx, err := tracing.InitJaeger("client", "127.0.0.1:6831")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Jaeger: %v", err))
	}
	defer tracer.Shutdown(ctx)
	ts := otel.Tracer("cccccc")
	plugins.Add(plugin.NewClientTracingPlugin(ts))
	xclient.SetPlugins(plugins)
	return tracer, ctx
}
