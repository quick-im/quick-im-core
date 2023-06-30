package message

import (
	"fmt"

	"github.com/smallnest/rpcx/server"
)

type rpcxServer struct {
	ip   string
	port uint16
}

func NewServer(ip string, port uint16) *rpcxServer {
	return &rpcxServer{
		ip:   ip,
		port: port,
	}
}

func (s *rpcxServer) Start() error {
	ser := server.NewServer()
	// ser.RegisterFunctionName(SERVER_NAME, SERVICE_GENERATE_MESSAGE_ID, s.GenerateMessageID, "")
	return ser.Serve("tcp", fmt.Sprintf("%s:%d", s.ip, s.port))
}
