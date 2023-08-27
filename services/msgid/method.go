package msgid

import (
	"context"
	"net"

	"github.com/quick-im/quick-im-core/services/msgid/internal/logic"
	"github.com/smallnest/rpcx/server"
)

type GenerateMessageIDArgs struct {
	ConversationType uint64
	ConversationID   string
}

type GenerateMessageIDReply struct {
	MsgID string
}

type generateMessageIdFn func(ctx context.Context, args GenerateMessageIDArgs, reply *GenerateMessageIDReply) error

func (s *rpcxServer) GenerateMessageID(ctx context.Context) generateMessageIdFn {
	return func(ctx context.Context, args GenerateMessageIDArgs, reply *GenerateMessageIDReply) error {
		clientConn := ctx.Value(server.RemoteConnContextKey).(net.Conn)
		println(clientConn.RemoteAddr().String())
		reply.MsgID = logic.GenerateRongCloudMessageID(args.ConversationType, args.ConversationID)
		return nil
	}
}
