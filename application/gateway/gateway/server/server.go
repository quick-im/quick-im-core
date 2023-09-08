package server

import (
	"fmt"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/logger/innerzap"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"go.uber.org/zap/zapcore"
)

type Option func(*apiServer)

type apiServer struct {
	ip                 string
	port               uint16
	openTracing        bool
	serviceName        string
	trackAgentHostPort string
	useConsulRegistry  bool
	consulServers      []string
	logger             logger.Logger
}

func NewApiServer(opt ...Option) *apiServer {
	ser := &apiServer{
		consulServers: make([]string, 0),
		serviceName:   "Gateway",
	}
	for i := range opt {
		opt[i](ser)
	}
	if ser.logger == nil {
		ser.logger = innerzap.NewZapLoggerAdapter(
			innerzap.NewLoggerWithOpt(
				innerzap.WithLogLevel(zapcore.DebugLevel),
				innerzap.WithServiceName(ser.serviceName),
				innerzap.WithLogPath("logs"),
			).NewLogger(),
		)
	}
	return ser
}

func WithIp(ip string) Option {
	return func(rs *apiServer) {
		rs.ip = ip
	}
}

func WithPort(port uint16) Option {
	return func(rs *apiServer) {
		rs.port = port
	}
}

func WithOpenTracing(disable bool) Option {
	return func(rs *apiServer) {
		rs.openTracing = disable
	}
}

func WithJeagerServiceName(serviceName string) Option {
	return func(rs *apiServer) {
		rs.serviceName = serviceName
	}
}

func WithJeagerAgentHostPort(trackAgentHostPort string) Option {
	return func(rs *apiServer) {
		rs.trackAgentHostPort = trackAgentHostPort
	}
}

func WithUseConsulRegistry(useConsulRegistry bool) Option {
	return func(rs *apiServer) {
		rs.useConsulRegistry = useConsulRegistry
	}
}

func WithConsulServer(consulServer string) Option {
	return func(rs *apiServer) {
		if rs.consulServers == nil {
			rs.consulServers = make([]string, 0)
		}
		rs.consulServers = append(rs.consulServers, consulServer)
	}
}

func WithConsulServers(consulServers ...string) Option {
	return func(rs *apiServer) {
		rs.consulServers = consulServers
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(rs *apiServer) {
		rs.logger = logger
	}
}

func (r *apiServer) InitDepServices(serviceName string) *rpcx.RpcxClientWithOpt {
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
