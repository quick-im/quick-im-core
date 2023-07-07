package conversation

import (
	"context"
	"fmt"
	"time"

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
			s.logger.Error("CreateConvercation ConversationType Err", "CreateConvercationArgsType:", fmt.Sprintf("%d", args.ConversationType))
			return quickim_errors.ErrConversationTypeRange
		}
		if len(args.SessionList) < 1 {
			s.logger.Error("CreateConvercation ConversationNumberRange Err", "args:", fmt.Sprintf("%+v", args))
			return quickim_errors.ErrConversationNumberRange
		}
		reply.ConversationID = uuid.New().String()
		if err := dbObj.CreateConvercation(ctx, reply.ConversationID); err != nil {
			s.logger.Error("CreateConvercation To Db Err", "err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
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
			s.logger.Error("CreateConvercation SessionJoinsConvercationUseCopyFrom Err", "err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
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
			s.logger.Error("JoinConvercation ConversationNumberRange Err", "args:", fmt.Sprintf("%+v", args))
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
			s.logger.Error("JoinConvercation SessionJoinsConvercationUseCopyFrom Err", "err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
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
	return func(ctx context.Context, args GetJoinedConversationsArgs, reply *GetJoinedConversationsReply) error {
		list, err := dbObj.GetJoinedConversations(ctx, args.SessionId)
		if err != nil {
			r.logger.Error("GetJoinedConversations GetJoinedConversations Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
			return err
		}
		reply.Conversations = make([]string, len(list))
		for i := range list {
			reply.Conversations[i] = list[i].ConvercationID
		}
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

type checkJoinedConversationFn func(ctx context.Context, args CheckJoinedConversationArgs, reply *CheckJoinedConversationReply) error

func (r *rpcxServer) CheckJoinedConversation(ctx context.Context) checkJoinedConversationFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args CheckJoinedConversationArgs, reply *CheckJoinedConversationReply) error {
		reply.Joined = false
		n, err := dbObj.CheckJoinedonversation(ctx, db.CheckJoinedonversationParams{
			SessionID:      args.SessionId,
			ConvercationID: args.ConversationId,
		})
		if err != nil {
			r.logger.Error("CheckJoinedConversation CheckJoinedonversation Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
			return err
		}
		if n > 0 {
			reply.Joined = true
		}
		return nil
	}
}

// 移出会话
type KickoutForConversationArgs struct {
	SessionId      []string
	ConversationId string
}

type KickoutForConversationReply struct {
	Failed []string
}

type kickoutJoinedConversationFn func(ctx context.Context, args KickoutForConversationArgs, reply *KickoutForConversationReply) error

func (r *rpcxServer) KickoutForConversation(ctx context.Context) kickoutJoinedConversationFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args KickoutForConversationArgs, reply *KickoutForConversationReply) error {
		params := make([]db.KickoutForConversationParams, len(args.SessionId))
		for i := range args.SessionId {
			params[i].SessionID = args.SessionId[i]
			params[i].ConvercationID = args.ConversationId
		}
		dbObj.KickoutForConversation(ctx, params).Exec(func(i int, err error) {
			if err != nil {
				if reply.Failed == nil {
					reply.Failed = make([]string, 0)
				}
				reply.Failed = append(reply.Failed, args.SessionId[i])
				r.logger.Error("KickoutForConversation KickoutForConversation Err:", fmt.Sprintf("record: %d,arg: %+v", i, params[i]), " err:", err.Error())
			}
		})
		return nil
	}
}

// 获取会话信息
type GetConversationInfoArgs struct {
	ConversationId string
}

type GetConversationInfoReply struct {
	db.Conversation
}

type GetConversationInfoFn func(ctx context.Context, args GetConversationInfoArgs, reply *GetConversationInfoReply) error

func (r *rpcxServer) GetConversationInfo(ctx context.Context) GetConversationInfoFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args GetConversationInfoArgs, reply *GetConversationInfoReply) error {
		info, err := dbObj.GetConversationInfo(ctx, args.ConversationId)
		if err != nil {
			r.logger.Error("GetConversationInfo GetConversationInfo Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
			return err
		}
		reply.Conversation = info
		return nil
	}
}

// 删除会话
type SetDeleteConversationArgs struct {
	ConversationId []string
}

type SetDeleteConversationReply struct {
	Failed []string
}

type SetDeleteConversationFn func(ctx context.Context, args SetDeleteConversationArgs, reply *SetDeleteConversationReply) error

func (r *rpcxServer) SetDeleteConversation(ctx context.Context) SetDeleteConversationFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args SetDeleteConversationArgs, reply *SetDeleteConversationReply) error {
		dbObj.DeleteConversations(ctx, args.ConversationId).Exec(func(i int, err error) {
			if err != nil {
				if reply.Failed == nil {
					reply.Failed = make([]string, 0)
				}
				reply.Failed = append(reply.Failed, args.ConversationId[i])
				r.logger.Error("SetDeleteConversation DeleteConversations Err:", fmt.Sprintf("record: %d,arg: %+v", i, args.ConversationId[i]), " err:", err.Error())
			}
		})
		return nil
	}
}

// 设置归档会话
type SetArchiveConversationsArgs struct {
	ConversationId []string
	IsArchive      bool
}

type SetArchiveConversationsReply struct {
	Failed []string
}

type SetArchiveConversationsFn func(ctx context.Context, args SetArchiveConversationsArgs, reply *SetArchiveConversationsReply) error

func (r *rpcxServer) SetArchiveConversations(ctx context.Context) SetArchiveConversationsFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args SetArchiveConversationsArgs, reply *SetArchiveConversationsReply) error {
		if args.IsArchive {
			dbObj.ArchiveConversations(ctx, args.ConversationId).Exec(func(i int, err error) {
				if err != nil {
					if reply.Failed == nil {
						reply.Failed = make([]string, 0)
					}
					reply.Failed = append(reply.Failed, args.ConversationId[i])
					r.logger.Error("SetArchiveConversations ArchiveConversations Err:", fmt.Sprintf("record: %d,arg: %+v", i, args.ConversationId[i]), " err:", err.Error())
				}
			})
		} else {
			dbObj.UnArchiveConversations(ctx, args.ConversationId).Exec(func(i int, err error) {
				if err != nil {
					if reply.Failed == nil {
						reply.Failed = make([]string, 0)
					}
					reply.Failed = append(reply.Failed, args.ConversationId[i])
					r.logger.Error("SetArchiveConversations UnArchiveConversations Err:", fmt.Sprintf("record: %d,arg: %+v", i, args.ConversationId[i]), " err:", err.Error())
				}
			})
		}
		return nil
	}
}

// 更新会话LastMsg
type UpdateConversationLastMsgArgs struct {
	ConversationId  string
	MsgId           string
	LastTime        *time.Time
	LastSendSession string
}

type UpdateConversationLastMsgReply struct {
}

type UpdateConversationLastMsgFn func(ctx context.Context, args UpdateConversationLastMsgArgs, reply *UpdateConversationLastMsgReply) error

func (r *rpcxServer) UpdateConversationLastMsg(ctx context.Context) UpdateConversationLastMsgFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args UpdateConversationLastMsgArgs, reply *UpdateConversationLastMsgReply) error {
		err := dbObj.UpdateConversationLastMsg(ctx, db.UpdateConversationLastMsgParams{
			LastSendTime:    args.LastTime,
			LastMsgID:       args.MsgId,
			LastSendSession: args.LastSendSession,
			ConversationID:  args.ConversationId,
		})
		if err != nil {
			r.logger.Error("UpdateConversationLastMsg UpdateConversationLastMsg Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
		}
		return err
	}
}
