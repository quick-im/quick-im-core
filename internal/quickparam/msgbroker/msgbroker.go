package msgbroker

import (
	"net"

	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

type BroadcastArgs = model.Msg

type BroadcastReply struct {
}

type RegisterSessionInfo struct {
	Platfotm  uint8
	Conn      net.Conn
	Uid       string
	SessionId string
}

type RegisterSessionReply struct {
}
