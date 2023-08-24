package msgbroker

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

// 用于广播消息到client
func (r *rpcxServer) listenMsg(ctx context.Context, nc *messaging.NatsWarp) {
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	var c codec.GobUtils[model.Msg]
	var msgData model.Msg
	sub, err := js.Subscribe(config.MqMsgBrokerSubject, func(msg *nats.Msg) {
		// 广播
		_ = c.Decode(msg.Data, &msgData)
		println(fmt.Sprintf("%#v", msgData))
		msg.Ack()
	}, nats.DeliverNew())
	if err != nil {
		r.logger.Warn("ListenMsg Err", err.Error())
	}
	defer sub.Unsubscribe()
	select {}
}
