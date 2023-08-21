package conversation

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/messaging"
)

// 用于更新会话状态
func (r *rpcxServer) listenMsg(nc *messaging.NatsWarp) {
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	sub, err := js.Subscribe(config.MqMsgBrokerSubject, func(msg *nats.Msg) {
		fmt.Println("m1", string(msg.Data))
		msg.Ack()
	}, nats.DeliverLast())
	if err != nil {
		r.logger.Warn("ListenMsg Err", err.Error())
	}
	defer sub.Unsubscribe()
}
