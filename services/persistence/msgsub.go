package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// 用于持久化存储消息
func (r *rpcxServer) listenMsg(ctx context.Context, nc *messaging.NatsWarp) {
	var rdb contant.RethinkDbCtxType
	rdb = helper.GetCtxValue[contant.RethinkDbCtxType](ctx, contant.CTX_RETHINK_DB_KEY, rdb)
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	var c codec.GobUtils[model.Msg]
	var msgData model.Msg
	sub, err := js.QueueSubscribe(config.MqMsgPersistenceGroup, "quickim-persistence", func(msg *nats.Msg) {
		if err := c.Decode(msg.Data, &msgData); err != nil {
			r.logger.Error("listenMsg SaveMsgToDb GobDecode Err", fmt.Sprintf("arg:%+v err:%v", msg.Data, err))
			_ = msg.Ack()
			return
		}
		result, err := rethinkdb.Table(config.RethinkMsgDb).
			Insert(msgData).
			RunWrite(rdb)
		if err != nil {
			_ = msg.Nak()
			r.logger.Error("SaveMsgToDb rethinkdb.Insert Err", fmt.Sprintf("result:%+v arg:%+v err:%v", result, msgData, err))
			return
		}
		_ = msg.Ack()
	}, nats.AckExplicit(), nats.AckWait(30*time.Second), nats.MaxDeliver(3))
	if err != nil {
		r.logger.Warn("ListenMsg Err", err.Error())
	}
	defer sub.Unsubscribe()
	select {}
}
