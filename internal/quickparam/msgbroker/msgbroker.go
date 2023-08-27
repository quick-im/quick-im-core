package msgbroker

import (
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

type BroadcastArgs = model.Msg

type BroadcastReply struct {
}

type RegisterSessionInfo struct {
	Platform  uint8
	Uid       string
	SessionId string
}

type RegisterSessionReply struct {
}
