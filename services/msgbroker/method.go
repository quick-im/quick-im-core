package msgbroker

import (
	"context"
	"fmt"
	"net"

	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"github.com/quick-im/quick-im-core/internal/quickparam/msgbroker"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/services/conversation"
	"github.com/smallnest/rpcx/server"
)

type BroadcastRecvFn func(context.Context, msgbroker.BroadcastArgs, *msgbroker.BroadcastReply) error

func (r *rpcxServer) BroadcastRecv(ctx context.Context) BroadcastRecvFn {
	var c codec.GobUtils[model.Msg]
	var conversationService *rpcx.RpcxClientWithOpt
	conversationService = helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, conversationService)
	return func(ctx context.Context, args msgbroker.BroadcastArgs, reply *msgbroker.BroadcastReply) error {
		getSessionsArgs := conversation.GetConversationSessionsArgs{
			ConversationId: args.ConversationID,
		}
		getSessionsReply := conversation.GetConversationSessionsReply{}
		err := conversationService.Call(ctx, conversation.SERVICE_GET_CONVERSATION_SSESSIONS, getSessionsArgs, &getSessionsReply)
		if err != nil {
			r.logger.Error("BroadcastRecv Call Service: conversationService Method: SERVICE_GET_CONVERSATION_SSESSIONS failed,", fmt.Sprintf("args: %#v,err: %v", args, err))
			return err
		}
		msg := model.Msg(args)
		data, err := c.Encode(msg)
		if err != nil {
			r.logger.Error("BroadcastRecv Msg EncodeErr:", fmt.Sprintf("args: %#v,err: %v", data, err))
			return err
		}
		for i := range getSessionsReply.Sessions {
			if c, exist := r.connList[getSessionsReply.Sessions[i]]; exist {
				if _, err := c.Conn.Write(data); err != nil {
					r.logger.Error("BroadcastRecv Send Msg To Session Err:", fmt.Sprintf("conm: %#v,err: %v", c, err))
				}
			}
		}
		return nil
	}
}

type RegisterSessionFn func(context.Context, msgbroker.RegisterSessionInfo, msgbroker.RegisterSessionReply) error

// 将接入层的Session注册到broker，用于消息分发，这里可以优化一下，在注册之前踢掉其他同Session同Platform的客户端
func (r *rpcxServer) RegisterSession(ctx context.Context) RegisterSessionFn {
	return func(ctx context.Context, args msgbroker.RegisterSessionInfo, reply msgbroker.RegisterSessionReply) error {
		clientConn := ctx.Value(server.RemoteConnContextKey).(net.Conn)
		args.Conn = clientConn
		r.connList[args.SessionId] = args
		return nil
	}
}
