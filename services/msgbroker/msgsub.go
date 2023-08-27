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
		r.clientList.lock.RLock()
		for i := range getSessionsReply.Sessions {
			if clientAddr, exist := r.clientList.sessonIndex[getSessionsReply.Sessions[i]]; exist {
				//TODO: 这里的data要包装一下，告诉client发送给具体的session
				if _, err := r.clientList.client[clientAddr].conn.Write(msg.Data); err != nil {
					r.logger.Error("BroadcastRecv Send Msg To Session Err:", fmt.Sprintf("session: %s, err: %v", getSessionsReply.Sessions[i], err))
				}
			}
		}
		r.clientList.lock.RUnlock()
		//
		_ = msg.Ack()
	}, nats.DeliverNew())
	if err != nil {
		r.logger.Warn("ListenMsg Err", err.Error())
	}
	defer sub.Unsubscribe()
	select {}
}
