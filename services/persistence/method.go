package persistence

import (
	"context"
	"fmt"
	"sort"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
)

// 持久化消息
type SaveMsgToDbArgs struct {
	Msgs []db.SaveMsgToDbParams
}

type SaveMsgToDbReply struct {
	Failed []string
}

type saveMsgToDbFn func(context.Context, SaveMsgToDbArgs, *SaveMsgToDbReply) error

func (r *rpcxServer) SaveMsgToDb(ctx context.Context) saveMsgToDbFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args SaveMsgToDbArgs, reply *SaveMsgToDbReply) error {
		dbObj.SaveMsgToDb(ctx, args.Msgs).Exec(func(i int, err error) {
			if reply.Failed == nil {
				reply.Failed = make([]string, 0)
			}
			reply.Failed = append(reply.Failed, args.Msgs[i].MsgID)
			r.logger.Error("SaveMsgToDb Err", fmt.Sprintf("record:%d arg:%+v err:%v", i, args.Msgs[i], err))
		})
		return nil
	}
}

// 根据范围会获取消息消息
type GetMsgFromDbInRangeArgs struct {
	ConvercationID string
	StartMsgId     string
	EndMsgId       string
	Sort           contant.Sort
}

type GetMsgFromDbInRangeReply struct {
	Msg []db.Message
}

type getMsgFromDbInRangeFn func(context.Context, GetMsgFromDbInRangeArgs, *GetMsgFromDbInRangeReply) error

func (r *rpcxServer) GetMsgFromDbInRange(ctx context.Context) getMsgFromDbInRangeFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args GetMsgFromDbInRangeArgs, reply *GetMsgFromDbInRangeReply) error {
		msg, err := dbObj.GetMsgFromDbInRange(ctx, db.GetMsgFromDbInRangeParams{
			ConvercationID: args.ConvercationID,
			StartMsgID:     args.StartMsgId,
			EndMsgID:       args.EndMsgId,
		})
		if err != nil {
			r.logger.Error("GetMsgFromDbInRange GetMsgFromDbInRange Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msg, func(i, j int) bool {
				return msg[i].MsgID < msg[j].MsgID
			})
		} else {
			sort.Slice(msg, func(i, j int) bool {
				return msg[i].MsgID > msg[j].MsgID
			})
		}
		reply.Msg = msg
		return nil
	}
}

// 获取会话最后30条消息
type GetLast30MsgFromDbArgs struct {
	ConvercationID string
	Sort           contant.Sort
}

type GetLast30MsgFromDbReply struct {
	Msg []db.Message
}

type getLast30MsgFromDbFn func(context.Context, GetLast30MsgFromDbArgs, *GetLast30MsgFromDbReply) error

func (r *rpcxServer) GetLast30MsgFromDb(ctx context.Context) getLast30MsgFromDbFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args GetLast30MsgFromDbArgs, reply *GetLast30MsgFromDbReply) error {
		msg, err := dbObj.GetLast30MsgFromDb(ctx, args.ConvercationID)
		if err != nil {
			r.logger.Error("GetLast30MsgFromDb GetLast30MsgFromDb Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msg, func(i, j int) bool {
				return msg[i].MsgID < msg[j].MsgID
			})
		} else {
			sort.Slice(msg, func(i, j int) bool {
				return msg[i].MsgID > msg[j].MsgID
			})
		}
		reply.Msg = msg
		return nil
	}
}

// 获取指定消息id之前30条消息
type GetThe30MsgBeforeTheIdArgs struct {
	ConvercationID string
	MsgId          string
	Sort           contant.Sort
}

type GetThe30MsgBeforeTheIdReply struct {
	Msg []db.Message
}

type getThe30MsgBeforeTheIdFn func(context.Context, GetThe30MsgBeforeTheIdArgs, *GetThe30MsgBeforeTheIdReply) error

func (r *rpcxServer) GetThe30MsgBeforeTheId(ctx context.Context) getThe30MsgBeforeTheIdFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args GetThe30MsgBeforeTheIdArgs, reply *GetThe30MsgBeforeTheIdReply) error {
		msg, err := dbObj.GetThe30MsgBeforeTheId(ctx, db.GetThe30MsgBeforeTheIdParams{
			ConvercationID: args.ConvercationID,
			MsgID:          args.MsgId,
		})
		if err != nil {
			r.logger.Error("GetThe30MsgBeforeTheId GetThe30MsgBeforeTheId Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msg, func(i, j int) bool {
				return msg[i].MsgID < msg[j].MsgID
			})
		} else {
			sort.Slice(msg, func(i, j int) bool {
				return msg[i].MsgID > msg[j].MsgID
			})
		}
		reply.Msg = msg
		return nil
	}
}

// 获取指定消息id之后30条消息
type GetThe30MsgAfterTheIdArgs struct {
	ConvercationID string
	MsgId          string
	Sort           contant.Sort
}

type GetThe30MsgAfterTheIdReply struct {
	Msg []db.Message
}

type getThe30MsgAfterTheIdFn func(context.Context, GetThe30MsgAfterTheIdArgs, *GetThe30MsgAfterTheIdReply) error

func (r *rpcxServer) GetThe30MsgAfterTheId(ctx context.Context) getThe30MsgAfterTheIdFn {
	ctxDb := ctx.Value(contant.CTX_POSTGRES_KEY).(contant.PgCtxType)
	dbObj := db.New(ctxDb)
	return func(ctx context.Context, args GetThe30MsgAfterTheIdArgs, reply *GetThe30MsgAfterTheIdReply) error {
		msg, err := dbObj.GetThe30MsgBeforeTheId(ctx, db.GetThe30MsgBeforeTheIdParams{
			ConvercationID: args.ConvercationID,
			MsgID:          args.MsgId,
		})
		if err != nil {
			r.logger.Error("GetThe30MsgAfterTheId GetThe30MsgBeforeTheId Err:", fmt.Sprintf("args: %+v, err: %v", args, err))
			return err
		}
		if args.Sort == contant.Asc {
			sort.Slice(msg, func(i, j int) bool {
				return msg[i].MsgID < msg[j].MsgID
			})
		} else {
			sort.Slice(msg, func(i, j int) bool {
				return msg[i].MsgID > msg[j].MsgID
			})
		}
		reply.Msg = msg
		return nil
	}
}
