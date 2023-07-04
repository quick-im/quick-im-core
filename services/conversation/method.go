package conversation

import (
	"context"

	"github.com/google/uuid"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/quickim_errors"
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
			return quickim_errors.ErrConversationTypeRange
		}
		if len(args.SessionList) < 1 {
			return quickim_errors.ErrConversationNumberRange
		}
		reply.ConversationID = uuid.New().String()
		if err := dbObj.CreateConvercation(ctx, reply.ConversationID); err != nil {
			return err
		}
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
		return nil
	}
}

// 加入会话
type JoinConvercationArgs = CreateConvercationArgs
type JoinConvercationReply = CreateConvercationReply
type JoinConvercationFn = createConvercationFn

func (s *rpcxServer) JoinConvercation(ctx context.Context) JoinConvercationFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args JoinConvercationArgs, reply *JoinConvercationReply) error {
		if len(args.SessionList) < 1 {
			return quickim_errors.ErrConversationNumberRange
		}
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
		return nil
	}
}
