package persistence

import (
	"context"
	"fmt"

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
