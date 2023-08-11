package persistence

import "github.com/quick-im/quick-im-core/internal/logger"

type Option func(*rpcxServer)

func WithIp(ip string) Option {
	return func(rs *rpcxServer) {
		rs.ip = ip
	}
}

func WithPort(port uint16) Option {
	return func(rs *rpcxServer) {
		rs.port = port
	}
}

func WithOpenTracing(disable bool) Option {
	return func(rs *rpcxServer) {
		rs.openTracing = disable
	}
}

func WithJeagerServiceName(serviceName string) Option {
	return func(rs *rpcxServer) {
		rs.serviceName = serviceName
	}
}

func WithJeagerAgentHostPort(trackAgentHostPort string) Option {
	return func(rs *rpcxServer) {
		rs.trackAgentHostPort = trackAgentHostPort
	}
}

func WithUseConsulRegistry(useConsulRegistry bool) Option {
	return func(rs *rpcxServer) {
		rs.useConsulRegistry = useConsulRegistry
	}
}

func WithConsulServer(consulServer string) Option {
	return func(rs *rpcxServer) {
		if rs.consulServers == nil {
			rs.consulServers = make([]string, 0)
		}
		rs.consulServers = append(rs.consulServers, consulServer)
	}
}

func WithConsulServers(consulServers ...string) Option {
	return func(rs *rpcxServer) {
		rs.consulServers = consulServers
	}
}

func WithNatsServer(natsServer string) Option {
	return func(rs *rpcxServer) {
		if rs.natsServers == nil {
			rs.natsServers = make([]string, 0)
		}
		rs.natsServers = append(rs.natsServers, natsServer)
	}
}

func WithNatsServers(natsServers ...string) Option {
	return func(rs *rpcxServer) {
		rs.natsServers = natsServers
	}
}

func WithNatsDisableJetstream() Option {
	return func(rs *rpcxServer) {
		rs.natsEnableJetstream = false
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(rs *rpcxServer) {
		rs.logger = logger
	}
}
