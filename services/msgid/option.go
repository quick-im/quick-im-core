package msgid

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
		rs.consulServers = append(rs.consulServers, consulServer)
	}
}

func WithConsulServers(consulServers ...string) Option {
	return func(rs *rpcxServer) {
		rs.consulServers = consulServers
	}
}
