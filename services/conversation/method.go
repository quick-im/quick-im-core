package conversation

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/quickerr"
)

// 创建会话
const (
	conversationTypeMax = 0xF
)

type CreateConversationArgs struct {
	ConversationType uint8
	SessionList      []string
}

type CreateConversationReply struct {
	ConversationID string
}

type createConversationFn func(ctx context.Context, args CreateConversationArgs, reply *CreateConversationReply) error

func (s *rpcxServer) CreateConversation(ctx context.Context) createConversationFn {
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	var cacheDb contant.CacheCtxType
	cacheDb = helper.GetCtxValue(ctx, contant.CTX_CACHE_DB_KEY, cacheDb)
	return func(ctx context.Context, args CreateConversationArgs, reply *CreateConversationReply) error {
		if args.ConversationType > conversationTypeMax {
			s.logger.Error("CreateConversation ConversationType Err", "CreateConversationArgsType:", fmt.Sprintf("%d", args.ConversationType))
			return quickerr.ErrConversationTypeRange
		}
		if len(args.SessionList) < 1 {
			s.logger.Error("CreateConversation ConversationNumberRange Err", "args:", fmt.Sprintf("%+v", args))
			return quickerr.ErrConversationNumberRange
		}
		reply.ConversationID = uuid.New().String()
		if err := dbObj.CreateConversation(ctx, reply.ConversationID); err != nil {
			s.logger.Error("CreateConversation dbObj.CreateConversation Err", "err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
			return err
		}
		sessions := make([]db.SessionJoinsConversationUseCopyFromParams, len(args.SessionList))
		for i := range args.SessionList {
			sessions[i] = db.SessionJoinsConversationUseCopyFromParams{
				SessionID:      args.SessionList[i],
				ConversationID: reply.ConversationID,
			}
		}
		if _, err := dbObj.SessionJoinsConversationUseCopyFrom(ctx, sessions); err != nil {
			s.logger.Error("CreateConversation dbObj.SessionJoinsConversationUseCopyFrom Err", "err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
			return err
		}
		if err := cacheDb.AddConverstaionSessions(reply.ConversationID, args.SessionList); err != nil {
			s.logger.Error("CreateConversation cacheDb.AddConverstaionSessions Err", "err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
		}
		return nil
	}
}

// 加入会话
type JoinConversationArgs = CreateConversationArgs
type JoinConversationReply = CreateConversationReply
type JoinConversationFn createConversationFn

func (s *rpcxServer) JoinConversation(ctx context.Context) JoinConversationFn {
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	var cacheDb contant.CacheCtxType
	cacheDb = helper.GetCtxValue(ctx, contant.CTX_CACHE_DB_KEY, cacheDb)
	return func(ctx context.Context, args JoinConversationArgs, reply *JoinConversationReply) error {
		if len(args.SessionList) < 1 {
			s.logger.Error("JoinConversation ConversationNumberRange Err", "args:", fmt.Sprintf("%+v", args))
			return quickerr.ErrConversationNumberRange
		}
		sessions := make([]db.SessionJoinsConversationUseCopyFromParams, len(args.SessionList))
		for i := range args.SessionList {
			sessions[i] = db.SessionJoinsConversationUseCopyFromParams{
				SessionID:      args.SessionList[i],
				ConversationID: reply.ConversationID,
			}
		}
		if _, err := dbObj.SessionJoinsConversationUseCopyFrom(ctx, sessions); err != nil {
			s.logger.Error("JoinConversation dbObj.SessionJoinsConversationUseCopyFrom Err", "err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
			return err
		}
		if err := cacheDb.AddConverstaionSessions(reply.ConversationID, args.SessionList); err != nil {
			s.logger.Error("JoinConversation cacheDb.AddConverstaionSessions Err", "err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
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
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args GetJoinedConversationsArgs, reply *GetJoinedConversationsReply) error {
		list, err := dbObj.GetJoinedConversations(ctx, args.SessionId)
		if err != nil {
			r.logger.Error("GetJoinedConversations dbObj.GetJoinedConversations Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
			return err
		}
		reply.Conversations = make([]string, len(list))
		for i := range list {
			reply.Conversations[i] = list[i].ConversationID
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
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	var cacheDb contant.CacheCtxType
	cacheDb = helper.GetCtxValue(ctx, contant.CTX_CACHE_DB_KEY, cacheDb)
	return func(ctx context.Context, args CheckJoinedConversationArgs, reply *CheckJoinedConversationReply) error {
		reply.Joined = false
		isExists, err := cacheDb.IsExistsInConversation(args.ConversationId, args.SessionId)
		if err != nil {
			r.logger.Error("CheckJoinedConversation cacheDb.IsExistsInConversation Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
		} else {
			reply.Joined = isExists
			return nil
		}
		n, err := dbObj.CheckJoinedonversation(ctx, db.CheckJoinedonversationParams{
			SessionID:      args.SessionId,
			ConversationID: args.ConversationId,
		})
		if err != nil {
			r.logger.Error("CheckJoinedConversation dbObj.CheckJoinedonversation Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
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
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	var cacheDb contant.CacheCtxType
	cacheDb = helper.GetCtxValue(ctx, contant.CTX_CACHE_DB_KEY, cacheDb)
	return func(ctx context.Context, args KickoutForConversationArgs, reply *KickoutForConversationReply) error {
		params := make([]db.KickoutForConversationParams, len(args.SessionId))
		for i := range args.SessionId {
			params[i].SessionID = args.SessionId[i]
			params[i].ConversationID = args.ConversationId
		}
		dbObj.KickoutForConversation(ctx, params).Exec(func(i int, err error) {
			if err != nil {
				if reply.Failed == nil {
					reply.Failed = make([]string, 0)
				}
				reply.Failed = append(reply.Failed, args.SessionId[i])
				r.logger.Error("KickoutForConversation dbObj.KickoutForConversation Err:", fmt.Sprintf("record: %d,arg: %+v", i, params[i]), " err:", err.Error())
			}
		})
		if err := cacheDb.DelConversationSession(args.ConversationId, args.SessionId); err != nil {
			r.logger.Error("KickoutForConversation cacheDb.DelConversationSession Err:", err.Error())
		}
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
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args GetConversationInfoArgs, reply *GetConversationInfoReply) error {
		info, err := dbObj.GetConversationInfo(ctx, args.ConversationId)
		if err != nil {
			r.logger.Error("GetConversationInfo dbObj.GetConversationInfo Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
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
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args SetDeleteConversationArgs, reply *SetDeleteConversationReply) error {
		dbObj.DeleteConversations(ctx, args.ConversationId).Exec(func(i int, err error) {
			if err != nil {
				if reply.Failed == nil {
					reply.Failed = make([]string, 0)
				}
				reply.Failed = append(reply.Failed, args.ConversationId[i])
				r.logger.Error("SetDeleteConversation dbObj.DeleteConversations Err:", fmt.Sprintf("record: %d,arg: %+v", i, args.ConversationId[i]), " err:", err.Error())
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
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args SetArchiveConversationsArgs, reply *SetArchiveConversationsReply) error {
		if args.IsArchive {
			dbObj.ArchiveConversations(ctx, args.ConversationId).Exec(func(i int, err error) {
				if err != nil {
					if reply.Failed == nil {
						reply.Failed = make([]string, 0)
					}
					reply.Failed = append(reply.Failed, args.ConversationId[i])
					r.logger.Error("SetArchiveConversations dbObj.ArchiveConversations Err:", fmt.Sprintf("record: %d,arg: %+v", i, args.ConversationId[i]), " err:", err.Error())
				}
			})
		} else {
			dbObj.UnArchiveConversations(ctx, args.ConversationId).Exec(func(i int, err error) {
				if err != nil {
					if reply.Failed == nil {
						reply.Failed = make([]string, 0)
					}
					reply.Failed = append(reply.Failed, args.ConversationId[i])
					r.logger.Error("SetArchiveConversations dbObj.UnArchiveConversations Err:", fmt.Sprintf("record: %d,arg: %+v", i, args.ConversationId[i]), " err:", err.Error())
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
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args UpdateConversationLastMsgArgs, reply *UpdateConversationLastMsgReply) error {
		err := dbObj.UpdateConversationLastMsg(ctx, db.UpdateConversationLastMsgParams{
			LastSendTime:    args.LastTime,
			LastMsgID:       args.MsgId,
			LastSendSession: args.LastSendSession,
			ConversationID:  args.ConversationId,
		})
		if err != nil {
			r.logger.Error("UpdateConversationLastMsg dbObj.UpdateConversationLastMsg Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
		}
		return err
	}
}

type GetConversationSessionsArgs struct {
	ConversationId string
}

type GetConversationSessionsReply struct {
	Sessions []string
}

type GetConversationSessionsFn func(ctx context.Context, args GetConversationSessionsArgs, reply *GetConversationSessionsReply) error

func (r *rpcxServer) GetConversationSessions(ctx context.Context) GetConversationSessionsFn {
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	var cacheDb contant.CacheCtxType
	cacheDb = helper.GetCtxValue(ctx, contant.CTX_CACHE_DB_KEY, cacheDb)
	return func(ctx context.Context, args GetConversationSessionsArgs, reply *GetConversationSessionsReply) error {
		sessions, err := cacheDb.GetConversationSessions(args.ConversationId)
		if err != nil {
			r.logger.Error("GetConversationSessions cacheDb.GetConversationSessions Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
			ids, err := dbObj.GetConversationsAllUsers(ctx)
			if err != nil {
				r.logger.Error("GetConversationSessions dbObj.GetConversationsAllUsers Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
				return err
			}
			for _, v := range ids {
				reply.Sessions = append(reply.Sessions, v.SessionID)
			}
		} else {
			reply.Sessions = sessions
		}
		if err := cacheDb.AddConverstaionSessions(args.ConversationId, reply.Sessions); err != nil {
			r.logger.Error("GetConversationSessions cacheDb.AddConverstaionSessions Err:", err.Error(), " arg:", fmt.Sprintf("%+v", args))
		}
		return nil
	}
}
