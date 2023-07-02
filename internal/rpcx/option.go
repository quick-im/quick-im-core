package rpcx

type Option func(*rpcxClient)

func SetOpenTracing(disable bool) Option {
	return func(rs *rpcxClient) {
		rs.openTracing = disable
	}
}

func SetJeagerServiceName(serviceName string) Option {
	return func(rs *rpcxClient) {
		rs.serviceName = serviceName
	}
}

func SetJeagerAgentHostPort(agentHostPort string) Option {
	return func(rs *rpcxClient) {
		rs.agentHostPort = agentHostPort
	}
}
