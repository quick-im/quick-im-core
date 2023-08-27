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

type broadcastRecvFn func(context.Context, msgbroker.BroadcastArgs, *msgbroker.BroadcastReply) error

func (r *rpcxServer) BroadcastRecv(ctx context.Context) broadcastRecvFn {
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
		r.connList.lock.RLock()
		for i := range getSessionsReply.Sessions {
			if c, exist := r.connList.connMap[getSessionsReply.Sessions[i]]; exist {
				for platform := range c.PlatformConn {
					if _, err := c.PlatformConn[platform].Write(data); err != nil {
						r.logger.Error("BroadcastRecv Send Msg To Session Err:", fmt.Sprintf("session: %s, platform: %d, err: %v", getSessionsReply.Sessions[i], platform, err))
					}
				}
			}
		}
		r.connList.lock.RUnlock()
		return nil
	}
}

type registerSessionFn func(context.Context, msgbroker.RegisterSessionInfo, *msgbroker.RegisterSessionReply) error

// 注册之前先发送广播踢掉同用户同平台的其他连接
func (r *rpcxServer) RegisterSession(ctx context.Context) registerSessionFn {
	// var c codec.GobUtils[model.Msg]
	var selfService *rpcx.RpcxClientWithOpt
	selfService = helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGBORKER, selfService)
	return func(ctx context.Context, args msgbroker.RegisterSessionInfo, reply *msgbroker.RegisterSessionReply) error {
		clientConn := ctx.Value(server.RemoteConnContextKey).(net.Conn)
		// args.PlatformConn = make(map[uint8]net.Conn)
		// 发送广播，踢掉其他重复的连接
		selfService.Broadcast(ctx, SERVICE_KICKOUT_DUPLICATE, args, reply)
		r.connList.lock.Lock()
		if info, ok := r.connList.connMap[args.SessionId]; ok {
			if info.PlatformConn == nil {
				info.PlatformConn = make(map[uint8]net.Conn)
			}
			// 将新连接的socket注册进来
			info.PlatformConn[args.Platform] = clientConn
		} else {
			r.connList.connMap[args.SessionId] = connInfo{
				PlatformConn: make(map[uint8]net.Conn),
				Uid:          args.Uid,
				SessionId:    args.SessionId,
			}
			r.connList.connMap[args.SessionId].PlatformConn[args.Platform] = clientConn
		}
		r.connList.lock.Unlock()
		return nil
	}
}

type kickoutDuplicateFn = registerSessionFn

func (r *rpcxServer) KickoutDuplicate(ctx context.Context) kickoutDuplicateFn {
	return func(ctx context.Context, rsi msgbroker.RegisterSessionInfo, rsr *msgbroker.RegisterSessionReply) error {
		if info, ok := r.connList.connMap[rsi.SessionId]; ok {
			r.connList.lock.RLock()
			if info.PlatformConn == nil {
				info.PlatformConn = make(map[uint8]net.Conn)
			} else if oldConn, ok := info.PlatformConn[rsi.Platform]; ok {
				// 如果该平台已经登陆，则踢掉，考虑先发送掉线信息
				oldConn.Close()
			}
			r.connList.lock.RUnlock()
		}
		return nil
	}
}
