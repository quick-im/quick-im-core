package conversation

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

// 用于更新会话状态,这里最好不要做成异步操作
func (r *rpcxServer) listenMsg(ctx context.Context, nc *messaging.NatsWarp) {
	var ctxDb contant.PgCtxType
	ctxDb = helper.GetCtxValue(ctx, contant.CTX_POSTGRES_KEY, ctxDb)
	dbObj := db.New(ctxDb)
	_ = dbObj
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	var c codec.GobUtils[model.Msg]
	var msgData model.Msg
	sub, err := js.Subscribe(config.MqMsgBrokerSubject, func(msg *nats.Msg) {
		fmt.Println("m1", string(msg.Data))
		err := c.Decode(msg.Data, &msgData)
		if err != nil {
			return
		}
		msg.Ack()
	}, nats.DeliverLast())
	if err != nil {
		r.logger.Warn("ListenMsg Err", err.Error())
	}
	defer sub.Unsubscribe()
	select {}
}
