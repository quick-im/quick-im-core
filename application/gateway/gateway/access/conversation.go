package access

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/quickerr"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/services/conversation"
)

type getConversationInfoArgs struct {
	ConversationID string `json:"conversation_id"`
}

func GetConversationInfo(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	_ = claims
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := getConversationInfoArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: GetConversationInfo ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		conversationInfoArgs := conversation.GetConversationInfoArgs{
			ConversationId: clientArgs.ConversationID,
		}
		conversationInfoReply := conversation.GetConversationInfoReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_GET_CONVERSATION_INFO, conversationInfoArgs, &conversationInfoReply); err != nil {
			log.Error("Gateway method: GetConversationInfo Fn conversationService.Call:conversation.SERVICE_GET_CONVERSATION_INFO ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(conversationInfoReply))
	}
}

func GetJoinedConversation(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		conversationGetJoinedArgs := conversation.GetJoinedConversationsArgs{
			SessionId: claims.Sid,
		}
		conversationGetJoinedReply := conversation.GetJoinedConversationsReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_GET_JOINED_CONVERSATIONS, conversationGetJoinedArgs, &conversationGetJoinedReply); err != nil {
			log.Error("Gateway method: GetJoinedConversation Fn conversationService.Call:conversation.SERVICE_GET_JOINED_CONVERSATIONS ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(conversationGetJoinedReply))
	}
}

type CheckJoinedConversationArgs struct {
	ConversationID string `json:"conversation_id"`
}

func CheckJoinedConversation(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := CheckJoinedConversationArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: CheckJoinedConversation ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		conversationCheckJoinedArgs := conversation.CheckJoinedConversationArgs{
			SessionId:      claims.Sid,
			ConversationId: clientArgs.ConversationID,
		}
		conversationCheckJoinedReply := conversation.CheckJoinedConversationReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_CHECK_JOINED_CONVERSATION, conversationCheckJoinedArgs, &conversationCheckJoinedReply); err != nil {
			log.Error("Gateway method: CheckJoinedConversation Fn conversationService.Call:conversation.SERVICE_CHECK_JOINED_CONVERSATION ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(conversationCheckJoinedReply))
	}
}

type createConversationArgs struct {
	ConversationType uint8    `json:"conversation_type"`
	Sessions         []string `json:"sessions"`
}

// todo: 创建会话成功的时候需要通知加入会话的用户
func CreateConversationInner(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := createConversationArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: CreateConversationInner ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		createConversationArgs := conversation.CreateConversationArgs{
			ConversationType: clientArgs.ConversationType,
			SessionList:      append(clientArgs.Sessions, claims.Sid),
		}
		createConversationReply := conversation.CreateConversationReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_CREATE_CONVERSATION, createConversationArgs, &createConversationReply); err != nil {
			log.Error("Gateway method: CreateConversationInner Fn conversationService.Call:conversation.SERVICE_CREATE_CONVERSATION ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(createConversationReply))
	}
}

type joinConversationArgs struct {
	ConversationID string   `json:"conversation_id"`
	Sessions       []string `json:"sessions"`
}

// 内部接口
func JoinConversationInner(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := joinConversationArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: JoinConversation ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		conversationJoinArgs := conversation.JoinConversationArgs{
			ConversationType: 0,
			SessionList:      []string{},
		}
		conversationJoinReply := conversation.JoinConversationReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_JOIN_CONVERSATION, conversationJoinArgs, &conversationJoinReply); err != nil {
			log.Error("Gateway method: JoinConversationInner Fn conversationService.Call:conversation.SERVICE_JOIN_CONVERSATION ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(conversationJoinReply))
	}
}
