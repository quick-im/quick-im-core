package msghub

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/services/msgbroker"
)

type SendMsgArgs struct {
	FromSession    string
	ConvercationId string
	MsgType        int32
	Content        []byte
	SendTime       time.Time
}

type SendMsgReply struct {
}

type sendMsgFn func(context.Context, SendMsgArgs, *SendMsgReply) error

func (r *rpcxServer) SendMsg(ctx context.Context) sendMsgFn {
	//TODO: 通过nats进行消息分发
	//注意：persistence和msgbroker多个相同示例要加入同一个组，防止消息重复处理
	nc, ok := ctx.Value(contant.CTX_NATS_KEY).(contant.NatsCtxType)
	if !ok {
		r.logger.Fatal("this service needs to rely on the nats service")
	}
	js, err := nc.JetStream()
	if err != nil {
		r.logger.Fatal("get nats jetstream err", fmt.Sprintf("%v", err))
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "MSG_STREAM",
		Subjects: []string{"stream.msg.>"},
	})
	if err != nil {
		r.logger.Fatal("add stream to nats jetstream err", fmt.Sprintf("%v", err))
	}
	msgbroker, err := rpcx.NewClient(
		rpcx.WithServerAddress(fmt.Sprintf("tcp@%s:%d", r.ip, r.port)),
		rpcx.WithUseConsulRegistry(r.useConsulRegistry),
		rpcx.WithConsulServers(r.consulServers...),
		rpcx.WithServiceName(msgbroker.SERVER_NAME),
		rpcx.WithOpenTracing(r.openTracing),
		rpcx.WithJeagerAgentHostPort(r.trackAgentHostPort),
	)
	if err != nil {
		r.logger.Fatal("init msgborker err", fmt.Sprintf("%v", err))
	}
	_ = msgbroker
	return func(ctx context.Context, args SendMsgArgs, reply *SendMsgReply) error {
		return nil
	}
}
