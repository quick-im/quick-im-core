package access

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/quickerr"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/services/conversation"
	"github.com/quick-im/quick-im-core/services/msghub"
	"github.com/quick-im/quick-im-core/services/msgid"
)

type sendMsgArgs struct {
	ConversationID string `json:"conversation_id"`
	Status         int32  `json:"status"`
	Type           int32  `json:"type"`
	Content        string `json:"content"`
}

func SendMsgHandler(ctx context.Context) http.HandlerFunc {
	// todo：这里的ctx断言每次调用都会执行，这里需要优化
	msghubService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGHUB, &rpcx.RpcxClientWithOpt{})
	msgidService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGID, &rpcx.RpcxClientWithOpt{})
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	// todo：这里发送消息的逻辑太过于冗长，后续版本需要优化
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := sendMsgArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: SendMsgHandler ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		conversationInfoArgs := conversation.GetConversationInfoArgs{
			ConversationId: clientArgs.ConversationID,
		}
		conversationInfoReply := conversation.GetConversationInfoReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_GET_CONVERSATION_INFO, conversationInfoArgs, &conversationInfoReply); err != nil {
			log.Error("Gateway method: SendMsgHandler Fn conversationService.Call:conversation.SERVICE_GET_CONVERSATION_INFO ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		msgIdArgs := msgid.GenerateMessageIDArgs{
			ConversationType: uint64(conversationInfoReply.ConversationType),
			ConversationID:   clientArgs.ConversationID,
		}
		msgIdReply := msgid.GenerateMessageIDReply{}
		if err := msgidService.Call(ctx, msgid.SERVICE_GENERATE_MESSAGE_ID, msgIdArgs, &msgIdReply); err != nil {
			log.Error("Gateway method: SendMsgHandler Fn msgidService.Call:msgid.SERVICE_GENERATE_MESSAGE_ID ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		sendMsgArgs := msghub.SendMsgArgs{
			MsgId:          msgIdReply.MsgID,
			FromSession:    claims.Sid,
			ConversationID: clientArgs.ConversationID,
			MsgType:        clientArgs.Type,
			Content:        []byte(clientArgs.Content),
			SendTime:       time.Now(),
		}
		sendMsgReply := msghub.SendMsgReply{}
		if err := msghubService.Call(ctx, msghub.SERVICE_SEND_MSG, sendMsgArgs, &sendMsgReply); err != nil {
			log.Error("Gateway method: SendMsgHandler Fn msghubService.Call:msghub.SERVICE_SEND_MSG ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		err := conversationService.Call(ctx, conversation.SERVICE_UPDATE_CONVERSATION_LASTMSG, conversation.UpdateConversationLastMsgArgs{
			ConversationId:  clientArgs.ConversationID,
			MsgId:           msgIdReply.MsgID,
			LastTime:        time.Now(),
			LastSendSession: claims.Sid,
		}, &conversation.UpdateConversationLastMsgReply{})
		if err != nil {
			log.Error("Gateway method: SendMsgHandler Fn conversationService.Call:conversation.SERVICE_UPDATE_CONVERSATION_LASTMSG ,err: ", err.Error())
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(msgIdReply))
	}
}
