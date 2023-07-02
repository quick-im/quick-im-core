package rpcx

type rpcxClient struct {
	openTracing   bool
	serviceName   string
	agentHostPort string
}

func NewClient(opts ...Option) *rpcxClient {
	ser := &rpcxClient{}
	for i := range opts {
		opts[i](ser)
	}
	return ser
}
