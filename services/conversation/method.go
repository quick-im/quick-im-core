package conversation

import (
	"context"

	"github.com/google/uuid"
	"github.com/quick-im/quick-im-core/internal/errors"
)

const (
	conversationTypeMax = 0xF
)

type CreateConvercationArgs struct {
	ConversationType uint8
	SessionList      []string
}

type CreateConvercationReply struct {
	ConversationID string
}

type createConvercationFn func(ctx context.Context, args CreateConvercationArgs, reply *CreateConvercationReply) error

func (s *rpcxServer) CreateConvercation(ctx context.Context) createConvercationFn {
	return func(ctx context.Context, args CreateConvercationArgs, reply *CreateConvercationReply) error {
		if args.ConversationType > conversationTypeMax {
			return errors.ErrConversationTypeRange
		}
		if len(args.SessionList) < 1 {
			return errors.ErrConversationNumberRange
		}
		reply.ConversationID = uuid.New().String()
		return nil
	}
}
