package rpcx

import (
	"context"
	"fmt"

	"github.com/quick-im/quick-im-core/internal/tracing"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/smallnest/rpcx/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type rpcxClient struct {
	openTracing   bool
	serviceName   string
	agentHostPort string
}

func NewClient(opts ...Option) *rpcxClient {
	c := &rpcxClient{}
	for i := range opts {
		opts[i](c)
	}
	// xclient := client.NewXClient(c.serviceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	// if c.openTracing {
	// 	tracer, ctx := c.addTrace(client)
	// 	_ = tracer
	// 	_ = ctx
	// }
	return c
}

func (s *rpcxClient) Call() error {
	return nil
}

func (s *rpcxClient) addTrace(xclient client.XClient) (*trace.TracerProvider, context.Context) {
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
