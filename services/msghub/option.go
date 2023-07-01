package msghub

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
