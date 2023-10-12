package msghub

import (
	"context"
	"fmt"

	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/services/msgbroker"
	"github.com/quick-im/quick-im-core/services/persistence"
)

// TODO：用户发送消息的前提是用户在会话中

type sendMsgFn func(context.Context, SendMsgArgs, *SendMsgReply) error

func (r *rpcxServer) SendMsg(ctx context.Context) sendMsgFn {
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
	gobc := codec.GobUtils[model.Msg]{}
	return func(ctx context.Context, args SendMsgArgs, reply *SendMsgReply) error {
		broadcastArgs := model.Msg{
			MsgId:          args.MsgId,
			ConversationID: args.ConversationID,
			FromSession:    args.FromSession,
			SendTime:       args.SendTime,
			Status:         0,
			Type:           args.MsgType,
			Content:        string(args.Content),
		}
		data, err := gobc.Encode(broadcastArgs)
		if err != nil {
			r.logger.Error("SendMsg Codec encoding of data failed. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
			return err
		}
		// 数据持久化
		// //TODO: 或许这里通过rpc同步调用持久化而非消息队列异步持久化更符合业务？
		_, err = js.Publish(config.MqMsgPersistenceGroup, data)
		if err != nil {
			r.logger.Error("SendMsg: push to nats MqMsgPersistenceGroup failed, started rpcx downgrade call. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
			args2 := persistence.SaveMsgToDbArgs{
				Msgs: []model.Msg{
					broadcastArgs,
				},
			}
			reply2 := persistence.SaveMsgToDbReply{}
			err = persistenceService.Call(ctx, persistence.SERVICE_SAVE_MSG_TO_DB, args2, &reply2)
			if err != nil {
				r.logger.Error("SendMsg: nats & rpcx call failed, failed to store message. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
				return err
			}
		}
		// TODO: 这里使用同步持久化的原因：防止消息推送比持久化先到达？Benchmark测试 同步比异步满3~4倍
		// args2 := persistence.SaveMsgToDbArgs{
		// 	Msgs: []model.Msg{
		// 		broadcastArgs,
		// 	},
		// }
		// reply2 := persistence.SaveMsgToDbReply{}
		// err = persistenceService.Call(ctx, persistence.SERVICE_SAVE_MSG_TO_DB, args2, &reply2)
		// if err != nil {
		// 	r.logger.Error("SendMsg: nats & rpcx call failed, failed to store message. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
		// 	return err
		// }
		// 消息广播给消息交付组件
		_, err = js.Publish(config.MqMsgBrokerSubject, data)
		if err != nil {
			r.logger.Error("SendMsg: push to nats MqMsgBrokerSubject failed, started rpcx downgrade call. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
			// 这里进行一下降级rpcx广播操作
			broadcastArgs := model.Msg{
				MsgId:          args.MsgId,
				ConversationID: args.ConversationID,
				FromSession:    args.FromSession,
				SendTime:       args.SendTime,
				Status:         0,
				Type:           args.MsgType,
				Content:        string(args.Content),
			}
			err := msgbrokerService.Broadcast(ctx, msgbroker.SERVICE_BROADCAST_RECV, broadcastArgs, &msgbroker.BroadcastReply{})
			if err != nil {
				r.logger.Error("SendMsg: nats & rpcx call failed, failed to send message. Err: ", fmt.Sprintf("arg:%+v err:%v", args, err))
				return err
			}
		}
		return nil
	}
}
