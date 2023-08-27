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
		// r.connList.lock.RLock()
		// for i := range getSessionsReply.Sessions {
		// 	if c, exist := r.connList.connMap[getSessionsReply.Sessions[i]]; exist {
		// 		for platform := range c.PlatformConn {
		// 			if _, err := c.PlatformConn[platform].Write(data); err != nil {
		// 				r.logger.Error("BroadcastRecv Send Msg To Session Err:", fmt.Sprintf("session: %s, platform: %d, err: %v", getSessionsReply.Sessions[i], platform, err))
		// 			}
		// 		}
		// 	}
		// }
		// r.connList.lock.RUnlock()
		// fix
		// 向该客户端连接的节点发送消息，再由节点发送给具体session
		r.clientList.lock.RLock()
		for i := range getSessionsReply.Sessions {
			if platforms, exist := r.clientList.sessonIndex[getSessionsReply.Sessions[i]]; exist {
				//TODO: 这里的data要包装一下，告诉client发送给具体的session
				for platform, gatewayUuid := range platforms {
					_ = platform
					if err := r.rpcxSer.SendMessage(r.clientList.client[gatewayUuid].conn, SERVER_NAME, SERVICE_BROADCAST_RECV, nil, data); err != nil {
						r.logger.Error("BroadcastRecv Send Msg To Session Err:", fmt.Sprintf("session: %s, err: %v", getSessionsReply.Sessions[i], err))
						return err
					}
				}
			}
		}
		r.clientList.lock.RUnlock()
		//
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
		_ = selfService.Broadcast(ctx, SERVICE_KICKOUT_DUPLICATE, args, reply)
		// r.connList.lock.Lock()
		// if info, ok := r.connList.connMap[args.SessionId]; ok {
		// 	if info.PlatformConn == nil {
		// 		info.PlatformConn = make(map[uint8]net.Conn)
		// 	}
		// 	// 将新连接的socket注册进来
		// 	info.PlatformConn[args.Platform] = clientConn
		// } else {
		// 	r.connList.connMap[args.SessionId] = connInfo{
		// 		PlatformConn: make(map[uint8]net.Conn),
		// 		Uid:          args.Uid,
		// 		SessionId:    args.SessionId,
		// 	}
		// 	r.connList.connMap[args.SessionId].PlatformConn[args.Platform] = clientConn
		// }
		// r.connList.lock.Unlock()
		// fix
		r.clientList.lock.Lock()
		if _, ok := r.clientList.client[args.GatewayUuid]; ok {
			// platforms := r.clientList.client[clientAddr].connMap[args.SessionId]
			// 如果session存在该节点则直接注册
			r.clientList.client[args.GatewayUuid].connMap[args.SessionId][args.Platform] = struct{}{}
		} else {
			// 如果不存在则重新注册
			r.clientList.client[args.GatewayUuid] = clientInfo{
				conn: clientConn,
				connMap: map[string]map[uint8]struct{}{
					args.SessionId: {
						args.Platform: struct{}{},
					},
				},
			}
			// 保存session && platform和网关的关联
			if r.clientList.sessonIndex[args.SessionId] == nil {
				// 初始化map
				r.clientList.sessonIndex[args.SessionId] = make(map[uint8]string)
			}
			r.clientList.sessonIndex[args.SessionId][args.Platform] = args.GatewayUuid
			// println("register ok")
			// println(r.clientList.sessonIndex[args.SessionId][args.Platform])
		}
		r.clientList.lock.Unlock()
		//
		return nil
	}
}

type kickoutDuplicateFn = registerSessionFn

func (r *rpcxServer) KickoutDuplicate(ctx context.Context) kickoutDuplicateFn {
	return func(ctx context.Context, rsi msgbroker.RegisterSessionInfo, rsr *msgbroker.RegisterSessionReply) error {
		// println("kictout")
		// if info, ok := r.connList.connMap[rsi.SessionId]; ok {
		// 	r.connList.lock.RLock()
		// 	if info.PlatformConn == nil {
		// 		info.PlatformConn = make(map[uint8]net.Conn)
		// 	} else if oldConn, ok := info.PlatformConn[rsi.Platform]; ok {
		// 		// 如果该平台已经登陆，则踢掉，考虑先发送掉线信息
		// 		oldConn.Close()
		// 	}
		// 	r.connList.lock.RUnlock()
		// }
		// fix
		r.clientList.lock.Lock()
		if platforms, ok := r.clientList.sessonIndex[rsi.SessionId]; ok {
			needDelete := false
			for platform, gatewayUuid := range platforms {
				if platform == rsi.Platform {
					// 如果该平台已登录
					// 删除索引
					if len(r.clientList.client[gatewayUuid].connMap[rsi.SessionId]) == 1 {
						// 如果只有一个待删除平台在这个节点，则直接删除session
						delete(r.clientList.client[gatewayUuid].connMap, rsi.SessionId)
					} else {
						// 如果还有其他平台在这个client节点，则只删除该平台
						delete(r.clientList.client[gatewayUuid].connMap[rsi.SessionId], rsi.Platform)
					}
					needDelete = true
					// println("这里踢出客户端")
					//TODO: 这里的data要包装一下，告诉client发送给具体的session
					_ = r.rpcxSer.SendMessage(r.clientList.client[gatewayUuid].conn, SERVER_NAME, SERVICE_KICKOUT_DUPLICATE, nil, []byte("kickout"))
					// 直接跳出处理，因为不该有第二个同用户的同平台在节点中，这是个bug
					break
				}
			}
			if needDelete {
				// 索引处理 {SessionIndex}
				// 如果该session只有一个platform在该msgbroker节点，则直接删除session索引,否则只删除对应platform
				if len(r.clientList.sessonIndex[rsi.SessionId]) == 1 {
					delete(r.clientList.sessonIndex, rsi.SessionId)
				} else {
					delete(r.clientList.sessonIndex[rsi.SessionId], rsi.Platform)
				}
			}
		}
		r.clientList.lock.Unlock()
		//
		return nil
	}
}
