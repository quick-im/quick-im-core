package msgbroker

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/services/conversation"
)

// 用于广播消息到client
func (r *rpcxServer) listenMsg(ctx context.Context, nc *messaging.NatsWarp) {
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	var c codec.GobUtils[model.Msg]
	var c2 codec.GobUtils[BroadcastMsgWarp]
	var msgData model.Msg
	var conversationService *rpcx.RpcxClientWithOpt
	conversationService = helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, conversationService)
	getSessionsArgs := conversation.GetConversationSessionsArgs{}
	getSessionsReply := conversation.GetConversationSessionsReply{}
	sub, err := js.Subscribe(config.MqMsgBrokerSubject, func(msg *nats.Msg) {
		// 广播
		if err := c.Decode(msg.Data, &msgData); err != nil {
			r.logger.Error("MsgBroker listenMsg Decode failed,", fmt.Sprintf(" args: %#v, err: %v", msg.Data, err))
			_ = msg.Ack()
			return
		}
		getSessionsArgs.ConversationId = msgData.ConversationID
		err := conversationService.Call(ctx, conversation.SERVICE_GET_CONVERSATION_SSESSIONS, getSessionsArgs, &getSessionsReply)
		if err != nil {
			r.logger.Error("MsgBroker Call Service: conversationService Method: SERVICE_GET_CONVERSATION_SSESSIONS failed,", fmt.Sprintf("args: %#v,err: %v", msgData, err))
			return
		}
		// r.connList.lock.RLock()
		// for i := range getSessionsReply.Sessions {
		// 	if c, exist := r.connList.connMap[getSessionsReply.Sessions[i]]; exist {
		// 		for platform := range c.PlatformConn {
		// 			if _, err := c.PlatformConn[platform].Write(msg.Data); err != nil {
		// 				r.logger.Error("MsgBroker Send Msg To Session Err:", fmt.Sprintf("session: %s, platform: %d, err: %v", getSessionsReply.Sessions[i], platform, err))
		// 			}
		// 		}
		// 	}
		// }
		// r.connList.lock.RUnlock()
		// fix
		// map[GatewayUuid]BroadcastMsgWarp
		recvSessions := BroadcastMsgWarp{
			Action:     SendMsg,
			MetaData:   msgData,
			ToSessions: []RecvSession{},
		}
		sendMaps := map[string][]RecvSession{}
		r.clientList.lock.RLock()
		defer r.clientList.lock.RUnlock()
		for i := range getSessionsReply.Sessions {
			if platforms, exist := r.clientList.sessonIndex[getSessionsReply.Sessions[i]]; exist {
				//TODO: 这里的data要包装一下，告诉client发送给具体的session
				for platform, gatewayUuid := range platforms {
					if sendMaps[gatewayUuid] == nil {
						sendMaps[gatewayUuid] = make([]RecvSession, 0)
					}
					sendMaps[gatewayUuid] = append(sendMaps[gatewayUuid], RecvSession{
						SessionId: getSessionsReply.Sessions[i],
						Platform:  platform,
					})
					// if err := r.rpcxSer.SendMessage(r.clientList.client[gatewayUuid].conn, SERVER_NAME, SERVICE_BROADCAST_RECV, nil, msg.Data); err != nil {
					// 	r.logger.Error("Msgbroker Send Msg To Session Err:", fmt.Sprintf("session: %s, platform: %d, err: %v", getSessionsReply.Sessions[i], platform, err))
					// }
				}
			}
		}
		// 将消息收集统一发送，减少数据包传输数量
		for gatewayUuid := range sendMaps {
			recvSessions.ToSessions = sendMaps[gatewayUuid]
			data, err := c2.Encode(recvSessions)
			if err != nil {
				r.logger.Error("MsgBroker listenMsg Encode failed,", fmt.Sprintf("args: %#v, err: %v", recvSessions, err))
				return
			}
			if err := r.rpcxSer.SendMessage(r.clientList.client[gatewayUuid].conn, SERVER_NAME, SERVICE_BROADCAST_RECV, nil, data); err != nil {
				r.logger.Error("Msgbroker Send Msg To Session Err:", fmt.Sprintf("gatewayUuid: %s, gatewayAddr: %s, err: %v", gatewayUuid, r.clientList.client[gatewayUuid].conn.RemoteAddr().String(), err))
			}
		}
		//
		_ = msg.Ack()
	}, nats.DeliverNew())
	if err != nil {
		r.logger.Warn("ListenMsg Err", err.Error())
	}
	defer sub.Unsubscribe()
	select {}
}
