package msgbroker

type Option func(*rpcxServer)

func SetOptIp(ip string) Option {
	return func(rs *rpcxServer) {
		rs.ip = ip
	}
}

func SetOptPort(port uint16) Option {
	return func(rs *rpcxServer) {
		rs.port = port
	}
}

func SetOpenTracing(disable bool) Option {
	return func(rs *rpcxServer) {
		rs.openTracing = disable
	}
}

func SetJeagerServiceName(disable bool) Option {
	return func(rs *rpcxServer) {
		rs.openTracing = disable
	}
}

func SetJeagerAgentHostPort(agentHostPort string) Option {
	return func(rs *rpcxServer) {
		rs.agentHostPort = agentHostPort
	}
}
