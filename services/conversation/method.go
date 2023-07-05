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
type JoinConvercationFn createConvercationFn

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

// 获取加入的会话
type GetJoinedConversationsArgs struct {
	SessionId string
}

type GetJoinedConversationsReply struct {
	Conversations []string
}

type getJoinedConversationsFn func(ctx context.Context, args GetJoinedConversationsArgs, reply *GetJoinedConversationsReply) error

func (r *rpcxServer) GetJoinedConversations(ctx context.Context) getJoinedConversationsFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	_ = dbObj
	return func(ctx context.Context, args GetJoinedConversationsArgs, reply *GetJoinedConversationsReply) error {
		return nil
	}
}

// 检测是否加入会话
type CheckJoinedConversationArgs struct {
	SessionId      string
	ConversationId string
}

type CheckJoinedConversationReply struct {
	Joined bool
}

type checkJoinedConversationsFn func(ctx context.Context, args CheckJoinedConversationArgs, reply *CheckJoinedConversationReply) error

func (r *rpcxServer) CheckJoinedConversations(ctx context.Context) checkJoinedConversationsFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	_ = dbObj
	return func(ctx context.Context, args CheckJoinedConversationArgs, reply *CheckJoinedConversationReply) error {
		return nil
	}
}

// 移出会话
type KickoutForConversationArgs struct {
	SessionId      []string
	ConversationId string
}

type KickoutForConversationReply struct {
	Kickouted bool
}

type kickoutJoinedConversationsFn func(ctx context.Context, args KickoutForConversationArgs, reply *KickoutForConversationReply) error

func (r *rpcxServer) KickoutForConversations(ctx context.Context) kickoutJoinedConversationsFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	_ = dbObj
	return func(ctx context.Context, args KickoutForConversationArgs, reply *KickoutForConversationReply) error {
		return nil
	}
}

// 获取会话信息
type GetConversationInfoArgs struct {
	ConversationId string
}

type GetConversationInfoReply struct {
	LasMsgId         string
	ConversationType int32
	IsDelete         bool
	IsArchive        bool
}

type GetConversationInfoFn func(ctx context.Context, args KickoutForConversationArgs, reply *KickoutForConversationReply) error

func (r *rpcxServer) GetConversationInfo(ctx context.Context) GetConversationInfoFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	_ = dbObj
	return func(ctx context.Context, args KickoutForConversationArgs, reply *KickoutForConversationReply) error {
		return nil
	}
}
