package persistence

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/messaging"
)

// 用于持久化存储消息
func (r *rpcxServer) listenMsg(nc *messaging.NatsWarp) {
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	sub, err := js.QueueSubscribe(config.MqMsgPersistenceGroup, "quickim-persistence", func(msg *nats.Msg) {
		fmt.Println("g1", string(msg.Data))
		msg.Ack()
	}, nats.AckExplicit(), nats.AckWait(30*time.Second))
	if err != nil {
		r.logger.Warn("ListenMsg Err", err.Error())
	}
	defer sub.Unsubscribe()
}
