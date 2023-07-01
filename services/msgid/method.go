package msgid

import (
	"context"

	"github.com/quick-im/quick-im-core/services/msgid/internal/logic"
)

type GenerateMessageIDArgs struct {
	ConversationType uint64
	ConversationID   string
}

type GenerateMessageIDReply struct {
	MsgID string
}

func (s *rpcxServer) GenerateMessageID(ctx context.Context, args GenerateMessageIDArgs, reply *GenerateMessageIDReply) error {
	reply.MsgID = logic.GenerateRongCloudMessageID(args.ConversationType, args.ConversationID)
	return nil
}
