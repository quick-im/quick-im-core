package msghub

import (
	"context"
	"fmt"
	"time"

	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/services/persistence"
)

type SendMsgArgs struct {
	MsgId          string
	FromSession    int32
	ConversationID string
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
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	var persistenceService *rpcx.RpcxClientWithOpt
	persistenceService = helper.GetCtxValue(ctx, contant.CTX_SERVICE_PERSISTENCE, persistenceService)
	var msgbrokerService *rpcx.RpcxClientWithOpt
	msgbrokerService = helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGBORKER, msgbrokerService)
	gobc := codec.GobUtils[SendMsgArgs]{}
	return func(ctx context.Context, args SendMsgArgs, reply *SendMsgReply) error {
		data, err := gobc.Encode(args)
		if err != nil {
			r.logger.Error("SendMsg Codec encoding of data failed. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
			return err
		}
		// 数据持久化
		_, err = js.Publish(config.MqMsgPersistenceGroup, data)
		if err != nil {
			r.logger.Error("SendMsg: push to nats MqMsgPersistenceGroup failed, started rpcx downgrade call. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
			reply := &persistence.SaveMsgToDbReply{}
			err = persistenceService.Call(context.Background(), persistence.SERVICE_SAVE_MSG_TO_DB, model.Msg{
				MsgId:          args.MsgId,
				ConversationID: args.ConversationID,
				FromSession:    args.FromSession,
				SendTime:       args.SendTime,
				Status:         0,
				Type:           args.MsgType,
				Content:        string(args.Content),
			}, reply)
			if err != nil {
				r.logger.Error("SendMsg: nats & rpcx call failed, failed to store message. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
				return err
			}
			// 消息广播给消息交付组件
			_, err = js.Publish(config.MqMsgBrokerSubject, data)
			if err != nil {
				r.logger.Error("SendMsg: push to nats MqMsgBrokerSubject failed, started rpcx downgrade call. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
				// 这里进行一下降级rpcx广播操作
				err := msgbrokerService.Broadcast(context.Background(), "method", "args", "reply")
				if err != nil {
					r.logger.Error("SendMsg: nats & rpcx call failed, failed to send message. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
					return err
				}
			}
		}
		return nil
	}
}
