package rpcx

import (
	"context"
	"errors"
	"fmt"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/quickerr"
	"github.com/quick-im/quick-im-core/internal/tracing"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	cclient "github.com/rpcxio/rpcx-consul/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

type RpcxClientWithOpt struct {
	openTracing              bool
	useConsulRegistry        bool
	basePath                 string
	serviceName              string
	clientName               string
	serverAddress            string
	consulServers            []string
	trackJaegeragentHostPort string
	tracePtr                 *trace.TracerProvider
	xclientPool              *client.XClientPool
	ctx                      context.Context
}

type metaDataWarp struct {
	key string
	val string
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

func WithClientName(clientName string) rpcxOption {
	return func(rs *RpcxClientWithOpt) {
		rs.clientName = clientName
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

func WithBasePath(basePath string) rpcxOption {
	return func(rco *RpcxClientWithOpt) {
		rco.basePath = basePath
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

func WithMetaData(key, val string) metaDataWarp {
	return metaDataWarp{
		key: key,
		val: val,
	}
}

func NewClient(opts ...rpcxOption) (*RpcxClientWithOpt, error) {
	c := &RpcxClientWithOpt{
		consulServers: make([]string, 0),
		basePath:      config.ServerPrefix,
	}
	for i := range opts {
		opts[i](c)
	}
	var cliDiscovery client.ServiceDiscovery
	var err error
	if len(c.consulServers) != 0 {
		cliDiscovery, err = cclient.NewConsulDiscovery(config.ServerPrefix, c.serviceName, c.consulServers, nil)
		if err != nil {
			return nil, err
		}
	} else {
		cliDiscovery, err = client.NewPeer2PeerDiscovery(c.serverAddress, "")
		if err != nil {
			return nil, err
		}
	}
	xclients := client.NewXClientPool(10, c.serviceName, client.Failtry, client.RandomSelect, cliDiscovery, client.DefaultOption)
	c.xclientPool = xclients
	return c, nil
}

func (c *RpcxClientWithOpt) GetOnce() (client.XClient, error) {
	var cliDiscovery client.ServiceDiscovery
	var err error
	if len(c.consulServers) != 0 {
		cliDiscovery, err = cclient.NewConsulDiscovery(config.ServerPrefix, c.serviceName, c.consulServers, nil)
		if err != nil {
			return nil, err
		}
	} else {
		cliDiscovery, err = client.NewPeer2PeerDiscovery(c.serverAddress, "")
		if err != nil {
			return nil, err
		}
	}
	xclient := client.NewXClient(c.serviceName, client.Failtry, client.RandomSelect, cliDiscovery, client.DefaultOption)
	if c.openTracing {
		c.addTrace(xclient)
	}
	return xclient, nil
}

func (s *RpcxClientWithOpt) Close() error {
	s.xclientPool.Close()
	return nil
}

func (s *RpcxClientWithOpt) ShutdownTrace() error {
	if s.tracePtr == nil {
		return quickerr.ErrTraceClosed
	}
	return s.tracePtr.Shutdown(s.ctx)
}

func (s *RpcxClientWithOpt) CloseAndShutdownTrace() error {
	var err error
	s.xclientPool.Close()
	if s.tracePtr != nil {
		err2 := s.tracePtr.Shutdown(s.ctx)
		if err2 != nil {
			err = errors.Join(err, err2)
		}
	}
	return err
}

func (s *RpcxClientWithOpt) Call(ctx context.Context, serviceMethod string, arg interface{}, replay interface{}, metadata ...metaDataWarp) error {
	if s.ctx == nil {
		s.ctx = ctx
	}
	meta := make(map[string]string, len(metadata))
	for i := range metadata {
		meta[metadata[i].key] = metadata[i].val
	}
	ctxInner := context.WithValue(s.ctx, ctxInitInner("initCtx"), nil)
	ctxInner = context.WithValue(ctxInner, share.ReqMetaDataKey, meta)
	xclient := s.xclientPool.Get()
	if s.openTracing {
		s.addTrace(xclient)
	}
	return xclient.Call(ctxInner, serviceMethod, arg, replay)
}

func (s *RpcxClientWithOpt) Broadcast(ctx context.Context, serviceMethod string, arg interface{}, replay interface{}, metadata ...metaDataWarp) error {
	if s.ctx == nil {
		s.ctx = ctx
	}
	meta := make(map[string]string, len(metadata))
	for i := range metadata {
		meta[metadata[i].key] = meta[metadata[i].val]
	}
	ctxInner := context.WithValue(s.ctx, ctxInitInner("initCtx"), nil)
	ctxInner = context.WithValue(ctxInner, share.ReqMetaDataKey, meta)
	xclient := s.xclientPool.Get()
	if s.openTracing {
		s.addTrace(xclient)
	}
	return xclient.Broadcast(ctxInner, serviceMethod, arg, replay)
}

func (s *RpcxClientWithOpt) addTrace(xclient client.XClient) {
	// 添加 Jaeger 拦截器
	plugins := client.NewPluginContainer()
	if s.tracePtr == nil {
		var err error
		s.tracePtr, s.ctx, err = tracing.InitJaeger(s.clientName, s.trackJaegeragentHostPort)
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize Jaeger: %v", err))
		}
	}
	ts := otel.Tracer(s.clientName)
	plugins.Add(plugin.NewClientTracingPlugin(ts))
	xclient.SetPlugins(plugins)
	// return tracer, ctx
}
