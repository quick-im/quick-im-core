package msgbroker

import (
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

type BroadcastArgs = model.Msg

type BroadcastReply struct {
}

type RegisterSessionInfo struct {
	Platform    uint8
	GatewayUuid string
	SessionId   string
}

type RegisterSessionReply struct {
}

type Action uint8

const (
	SendMsg = 1
	Kickout = 2
)

type BroadcastMsgWarp struct {
	Action     Action
	MetaData   model.Msg
	ToSessions []RecvSession
}

type RecvSession struct {
	SessionId string
	Platform  uint8
}
