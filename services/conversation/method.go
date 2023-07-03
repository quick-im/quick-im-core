package conversation

import (
	"context"

	"github.com/google/uuid"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/errors"
)

// 创建会话
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
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args CreateConvercationArgs, reply *CreateConvercationReply) error {
		if args.ConversationType > conversationTypeMax {
			return errors.ErrConversationTypeRange
		}
		if len(args.SessionList) < 1 {
			return errors.ErrConversationNumberRange
		}
		reply.ConversationID = uuid.New().String()
		if err := dbObj.CreateConvercation(ctx, reply.ConversationID); err != nil {
			return err
		}
		if len(args.SessionList) > 0 {
			sessions := make([]db.SessionJoinsConvercationUseCopyFromParams, len(args.SessionList))
			for i := range args.SessionList {
				sessions[i] = db.SessionJoinsConvercationUseCopyFromParams{
					SessionID:      args.SessionList[i],
					ConvercationID: reply.ConversationID,
				}
			}
			if _, err := dbObj.SessionJoinsConvercationUseCopyFrom(ctx, sessions); err != nil {
				return err
			}
		}
		return nil
	}
}
