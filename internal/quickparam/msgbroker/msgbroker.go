package msgbroker

import (
	"net"

	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

type BroadcastArgs = model.Msg

type BroadcastReply struct {
}

type RegisterSessionInfo struct {
	PlatformConn map[uint8]net.Conn
	Platform     uint8
	Uid          string
	SessionId    string
}

type RegisterSessionReply struct {
}
