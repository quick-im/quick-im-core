package msghub

import (
	"context"
	"time"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
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
	var nc contant.NatsCtxType
	nc = helper.GetCtxValue(ctx, contant.CTX_NATS_KEY, nc)
	_ = nc
	return func(ctx context.Context, args SendMsgArgs, reply *SendMsgReply) error {
		return nil
	}
}
