package persistence

import (
	"context"
	"fmt"
	"sort"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// 持久化消息
type SaveMsgToDbArgs struct {
	Msgs []model.Msg
}

type SaveMsgToDbReply struct {
	Failed []string
}

type saveMsgToDbFn func(context.Context, SaveMsgToDbArgs, *SaveMsgToDbReply) error

func (r *rpcxServer) SaveMsgToDb(ctx context.Context) saveMsgToDbFn {
	var rdb contant.RethinkDbCtxType
	rdb = helper.GetCtxValue[contant.RethinkDbCtxType](ctx, contant.CTX_RETHINK_DB_KEY, rdb)
	return func(ctx context.Context, args SaveMsgToDbArgs, reply *SaveMsgToDbReply) error {
		result, err := rethinkdb.Table(config.RethinkMsgDb).
			Insert(args.Msgs).
			RunWrite(rdb)
		if err != nil {
			r.logger.Error("SaveMsgToDb rethinkdb.Insert Err", fmt.Sprintf("result:%+v arg:%+v err:%v", result, args, err))
		}
		return err
	}
}

// 根据范围会获取消息消息
type GetMsgFromDbInRangeArgs struct {
	ConversationID string
	StartMsgId     string
	EndMsgId       string
	Sort           contant.Sort
}

type GetMsgFromDbInRangeReply struct {
	Msg []model.Msg
}

type getMsgFromDbInRangeFn func(context.Context, GetMsgFromDbInRangeArgs, *GetMsgFromDbInRangeReply) error

func (r *rpcxServer) GetMsgFromDbInRange(ctx context.Context) getMsgFromDbInRangeFn {
	var rdb contant.RethinkDbCtxType
	rdb = helper.GetCtxValue[contant.RethinkDbCtxType](ctx, contant.CTX_RETHINK_DB_KEY, rdb)
	return func(ctx context.Context, args GetMsgFromDbInRangeArgs, reply *GetMsgFromDbInRangeReply) error {
		var msgs []model.Msg
		rows, err := rethinkdb.Table(config.RethinkMsgDb).
			Filter(
				rethinkdb.Row.Field("conversation_id").Eq(args.ConversationID).
					And(rethinkdb.Row.Field("msg_id").Ge(args.StartMsgId)).
					And(rethinkdb.Row.Field("msg_id").Le(args.EndMsgId)),
			).
			Run(rdb)
		if err != nil {
			r.logger.Error("GetMsgFromDbInRange rethinkdb.GetMsgFromDbInRange Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		defer rows.Close()
		if err := rows.All(&msgs); err != nil {
			r.logger.Error("GetMsgFromDbInRange rethinkdb.Bind To Struct Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId < msgs[j].MsgId
			})
		} else {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId > msgs[j].MsgId
			})
		}
		reply.Msg = msgs
		return nil
	}
}

// 获取会话最后30条消息
type GetLast30MsgFromDbArgs struct {
	ConversationID string
	Sort           contant.Sort
}

type GetLast30MsgFromDbReply struct {
	Msg []model.Msg
}

type getLast30MsgFromDbFn func(context.Context, GetLast30MsgFromDbArgs, *GetLast30MsgFromDbReply) error

func (r *rpcxServer) GetLast30MsgFromDb(ctx context.Context) getLast30MsgFromDbFn {
	var rdb contant.RethinkDbCtxType
	rdb = helper.GetCtxValue[contant.RethinkDbCtxType](ctx, contant.CTX_RETHINK_DB_KEY, rdb)
	return func(ctx context.Context, args GetLast30MsgFromDbArgs, reply *GetLast30MsgFromDbReply) error {
		var msgs []model.Msg
		rows, err := rethinkdb.Table(config.RethinkMsgDb).
			Filter(
				rethinkdb.Row.Field("conversation_id").Eq(args.ConversationID),
			).
			OrderBy(rethinkdb.Desc(rethinkdb.Row.Field("msg_id"))).
			Limit(30).
			Run(rdb)
		if err != nil {
			r.logger.Error("GetLast30MsgFromDb rethinkdb.GetLast30MsgFromDb Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		defer rows.Close()
		if err := rows.All(&msgs); err != nil {
			r.logger.Error("GetLast30MsgFromDb rethinkdb.Bind To Struct Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId < msgs[j].MsgId
			})
		} else {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId > msgs[j].MsgId
			})
		}
		reply.Msg = msgs
		return nil
	}
}

// 获取指定消息id之前30条消息
type GetThe30MsgBeforeTheIdArgs struct {
	ConversationID string
	MsgId          string
	Sort           contant.Sort
}

type GetThe30MsgBeforeTheIdReply struct {
	Msg []model.Msg
}

type getThe30MsgBeforeTheIdFn func(context.Context, GetThe30MsgBeforeTheIdArgs, *GetThe30MsgBeforeTheIdReply) error

func (r *rpcxServer) GetThe30MsgBeforeTheId(ctx context.Context) getThe30MsgBeforeTheIdFn {
	var rdb contant.RethinkDbCtxType
	rdb = helper.GetCtxValue[contant.RethinkDbCtxType](ctx, contant.CTX_RETHINK_DB_KEY, rdb)
	return func(ctx context.Context, args GetThe30MsgBeforeTheIdArgs, reply *GetThe30MsgBeforeTheIdReply) error {
		var msgs []model.Msg
		rows, err := rethinkdb.Table(config.RethinkMsgDb).
			Filter(
				rethinkdb.Row.Field("conversation_id").Eq(args.ConversationID).
					And(rethinkdb.Row.Field("msg_id").Le(args.MsgId)),
			).
			OrderBy(rethinkdb.Desc(rethinkdb.Row.Field("msg_id"))).
			Limit(30).
			Run(rdb)
		if err != nil {
			r.logger.Error("GetThe30MsgBeforeTheId rethinkdb.GetThe30MsgBeforeTheId Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		defer rows.Close()
		if err := rows.All(&msgs); err != nil {
			r.logger.Error("GetThe30MsgBeforeTheId rethinkdb.Bind To Struct Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId < msgs[j].MsgId
			})
		} else {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId > msgs[j].MsgId
			})
		}
		reply.Msg = msgs
		return nil
	}
}

// 获取指定消息id之后30条消息
type GetThe30MsgAfterTheIdArgs struct {
	ConversationID string
	MsgId          string
	Sort           contant.Sort
}

type GetThe30MsgAfterTheIdReply struct {
	Msg []model.Msg
}

type getThe30MsgAfterTheIdFn func(context.Context, GetThe30MsgAfterTheIdArgs, *GetThe30MsgAfterTheIdReply) error

func (r *rpcxServer) GetThe30MsgAfterTheId(ctx context.Context) getThe30MsgAfterTheIdFn {
	var rdb contant.RethinkDbCtxType
	rdb = helper.GetCtxValue[contant.RethinkDbCtxType](ctx, contant.CTX_RETHINK_DB_KEY, rdb)
	return func(ctx context.Context, args GetThe30MsgAfterTheIdArgs, reply *GetThe30MsgAfterTheIdReply) error {
		var msgs []model.Msg
		rows, err := rethinkdb.Table(config.RethinkMsgDb).
			Filter(
				rethinkdb.Row.Field("conversation_id").Eq(args.ConversationID).
					And(rethinkdb.Row.Field("msg_id").Ge(args.MsgId)),
			).
			OrderBy(rethinkdb.Asc(rethinkdb.Row.Field("msg_id"))).
			Limit(30).
			Run(rdb)
		if err != nil {
			r.logger.Error("GetThe30MsgAfterTheId rethinkdb.GetThe30MsgBeforeTheId Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		defer rows.Close()
		if err := rows.All(&msgs); err != nil {
			r.logger.Error("GetThe30MsgAfterTheId rethinkdb.Bind To Struct Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId < msgs[j].MsgId
			})
		} else {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId > msgs[j].MsgId
			})
		}
		reply.Msg = msgs
		return nil
	}
}

// 获取会话最后一条消息
type GetLastOneMsgFromDbArgs struct {
	ConversationID string
	Sort           contant.Sort
}

type GetLastOneMsgFromDbReply struct {
	Msg []model.Msg
}

type getLastOneMsgFromDbFn func(context.Context, GetLastOneMsgFromDbArgs, *GetLastOneMsgFromDbReply) error

func (r *rpcxServer) GetLastOneMsgFromDb(ctx context.Context) getLastOneMsgFromDbFn {
	var rdb contant.RethinkDbCtxType
	rdb = helper.GetCtxValue[contant.RethinkDbCtxType](ctx, contant.CTX_RETHINK_DB_KEY, rdb)
	return func(ctx context.Context, args GetLastOneMsgFromDbArgs, reply *GetLastOneMsgFromDbReply) error {
		var msgs []model.Msg
		rows, err := rethinkdb.Table(config.RethinkMsgDb).
			Filter(
				rethinkdb.Row.Field("conversation_id").Eq(args.ConversationID),
			).
			OrderBy(rethinkdb.Desc(rethinkdb.Row.Field("msg_id"))).
			Limit(1).
			Run(rdb)
		if err != nil {
			r.logger.Error("GetLast30MsgFromDb rethinkdb.GetLast30MsgFromDb Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		defer rows.Close()
		if err := rows.All(&msgs); err != nil {
			r.logger.Error("GetLast30MsgFromDb rethinkdb.Bind To Struct Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId < msgs[j].MsgId
			})
		} else {
			sort.Slice(msgs, func(i, j int) bool {
				return msgs[i].MsgId > msgs[j].MsgId
			})
		}
		reply.Msg = msgs
		return nil
	}
}
