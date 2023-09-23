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

type joinConversationArgs struct {
	ConversationID string `json:"conversation_id"`
}

// 外部接口
func JoinConversation(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := joinConversationArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: JoinConversation ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		conversationJoinArgs := conversation.JoinConversationArgs{
			ConversationID: clientArgs.ConversationID,
			SessionList:    []string{claims.Sid},
		}
		conversationJoinReply := conversation.JoinConversationReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_JOIN_CONVERSATION, conversationJoinArgs, &conversationJoinReply); err != nil {
			log.Error("Gateway method: JoinConversation Fn conversationService.Call:conversation.SERVICE_JOIN_CONVERSATION ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(conversationJoinReply))
	}
}

type leaveConversationArgs struct {
	ConversationID string `json:"conversation_id"`
}

// 外部接口,离开会话
func LeaveConversation(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := leaveConversationArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: LeaveConversation ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		leaveConversationArgs := conversation.KickoutForConversationArgs{
			SessionId:      []string{claims.Sid},
			ConversationId: clientArgs.ConversationID,
		}
		leaveConversationReply := conversation.KickoutForConversationReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_KICKOUT_FOR_CONVERSATION, leaveConversationArgs, &leaveConversationReply); err != nil {
			log.Error("Gateway method: LeaveConversation Fn conversationService.Call:conversation.SERVICE_KICKOUT_FOR_CONVERSATION ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(leaveConversationReply))
	}
}

type createConversationArgs struct {
	ConversationType uint64   `json:"conversation_type"`
	Sessions         []string `json:"sessions"`
}

// todo: 创建会话成功的时候需要通知加入会话的用户
func CreateConversation(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	msghubService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGHUB, &rpcx.RpcxClientWithOpt{})
	msgidService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGID, &rpcx.RpcxClientWithOpt{})
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := createConversationArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: CreateConversation ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		createConversationArgs := conversation.CreateConversationArgs{
			ConversationType: clientArgs.ConversationType,
			SessionList:      append(clientArgs.Sessions, claims.Sid),
		}
		createConversationReply := conversation.CreateConversationReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_CREATE_CONVERSATION, createConversationArgs, &createConversationReply); err != nil {
			log.Error("Gateway method: CreateConversation Fn conversationService.Call:conversation.SERVICE_CREATE_CONVERSATION ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		msgIdArgs := msgid.GenerateMessageIDArgs{
			ConversationType: clientArgs.ConversationType,
			ConversationID:   createConversationReply.ConversationID,
		}
		msgIdReply := msgid.GenerateMessageIDReply{}
		if err := msgidService.Call(ctx, msgid.SERVICE_GENERATE_MESSAGE_ID, msgIdArgs, &msgIdReply); err != nil {
			log.Error("Gateway method: CreateConversation Fn msgidService.Call:msgid.SERVICE_GENERATE_MESSAGE_ID ,err: ", err.Error())
			// _ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			// return
		}
		sendMsgArgs := msghub.SendMsgArgs{
			MsgId:          msgIdReply.MsgID,
			FromSession:    claims.Sid,
			ConversationID: createConversationReply.ConversationID,
			MsgType:        0,
			Content:        []byte("您已加入对话 "),
			SendTime:       time.Now(),
		}
		sendMsgReply := msghub.SendMsgReply{}
		if err := msghubService.Call(ctx, msghub.SERVICE_SEND_MSG, sendMsgArgs, &sendMsgReply); err != nil {
			log.Error("Gateway method: CreateConversation Fn msghubService.Call:msghub.SERVICE_SEND_MSG ,err: ", err.Error())
			// _ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			// return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(createConversationReply))
	}
}

type createConversationInnerArgs struct {
	ConversationType uint64   `json:"conversation_type"`
	Sessions         []string `json:"sessions"`
}

// todo: 创建会话成功的时候需要通知加入会话的用户
func CreateConversationInner(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	msghubService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGHUB, &rpcx.RpcxClientWithOpt{})
	msgidService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGID, &rpcx.RpcxClientWithOpt{})
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	// claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := createConversationInnerArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: CreateConversationInner ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		createConversationArgs := conversation.CreateConversationArgs{
			ConversationType: clientArgs.ConversationType,
			SessionList:      clientArgs.Sessions,
		}
		createConversationReply := conversation.CreateConversationReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_CREATE_CONVERSATION, createConversationArgs, &createConversationReply); err != nil {
			log.Error("Gateway method: CreateConversationInner Fn conversationService.Call:conversation.SERVICE_CREATE_CONVERSATION ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		msgIdArgs := msgid.GenerateMessageIDArgs{
			ConversationType: clientArgs.ConversationType,
			ConversationID:   createConversationReply.ConversationID,
		}
		msgIdReply := msgid.GenerateMessageIDReply{}
		if err := msgidService.Call(ctx, msgid.SERVICE_GENERATE_MESSAGE_ID, msgIdArgs, &msgIdReply); err != nil {
			log.Error("Gateway method: CreateConversationInner Fn msgidService.Call:msgid.SERVICE_GENERATE_MESSAGE_ID ,err: ", err.Error())
			// _ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			// return
		}
		sendMsgArgs := msghub.SendMsgArgs{
			MsgId:          msgIdReply.MsgID,
			FromSession:    "",
			ConversationID: createConversationReply.ConversationID,
			MsgType:        0,
			Content:        []byte("您已加入对话 "),
			SendTime:       time.Now(),
		}
		sendMsgReply := msghub.SendMsgReply{}
		if err := msghubService.Call(ctx, msghub.SERVICE_SEND_MSG, sendMsgArgs, &sendMsgReply); err != nil {
			log.Error("Gateway method: CreateConversationInner Fn msghubService.Call:msghub.SERVICE_SEND_MSG ,err: ", err.Error())
			// _ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			// return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(createConversationReply))
	}
}

type joinConversationInnerArgs struct {
	ConversationID string   `json:"conversation_id"`
	Sessions       []string `json:"sessions"`
}

// 内部接口
func JoinConversationInner(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := joinConversationInnerArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: JoinConversation ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		conversationJoinArgs := conversation.JoinConversationArgs{
			ConversationID: clientArgs.ConversationID,
			SessionList:    clientArgs.Sessions,
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

type KickoutConversationArgs struct {
	ConversationID string   `json:"conversation_id"`
	Sessions       []string `json:"sessions"`
}

// 内部接口,踢出用户
func KickoutConversationInner(ctx context.Context) http.HandlerFunc {
	conversationService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_CONVERSATION, &rpcx.RpcxClientWithOpt{})
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := KickoutConversationArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: KickoutConversationInner ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		leaveConversationArgs := conversation.KickoutForConversationArgs{
			SessionId:      clientArgs.Sessions,
			ConversationId: clientArgs.ConversationID,
		}
		leaveConversationReply := conversation.KickoutForConversationReply{}
		if err := conversationService.Call(ctx, conversation.SERVICE_KICKOUT_FOR_CONVERSATION, leaveConversationArgs, &leaveConversationReply); err != nil {
			log.Error("Gateway method: KickoutConversationInner Fn conversationService.Call:conversation.SERVICE_KICKOUT_FOR_CONVERSATION ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(leaveConversationReply))
	}
}
