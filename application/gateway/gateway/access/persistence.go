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
	"github.com/quick-im/quick-im-core/services/persistence"
)

type getMsgFromDbInRangeArgs struct {
	ConversationID string `json:"conversation_id"`
	StartMsgId     string `json:"start_msg_id"`
	EndMsgId       string `json:"end_msg_id"`
	Desc           bool   `json:"desc"`
}

// 根据消息id范围获取消息
func GetMsgFromDbInRange(ctx context.Context) http.HandlerFunc {
	persistenceService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_PERSISTENCE, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	_ = claims
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := getMsgFromDbInRangeArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: GetMsgFromDbInRange ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		GetMsgFromDbInRangeArgs := persistence.GetMsgFromDbInRangeArgs{
			ConversationID: clientArgs.ConversationID,
			StartMsgId:     clientArgs.StartMsgId,
			EndMsgId:       clientArgs.EndMsgId,
			Sort:           contant.Sort(clientArgs.Desc),
		}
		GetMsgFromDbInRangeReply := persistence.GetMsgFromDbInRangeReply{}
		if err := persistenceService.Call(ctx, persistence.SERVICE_GET_MSG_FROM_DB_IN_RANGE, GetMsgFromDbInRangeArgs, &GetMsgFromDbInRangeReply); err != nil {
			log.Error("Gateway method: GetMsgFromDbInRange Fn persistenceService.Call:persistence.SERVICE_GET_MSG_FROM_DB_IN_RANGE ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(GetMsgFromDbInRangeReply))
	}
}

type getLast30MsgFromDbArgs struct {
	ConversationID string `json:"conversation_id"`
	Desc           bool   `json:"desc"`
}

// 获取会话最后30条消息
func GetLast30MsgFromDb(ctx context.Context) http.HandlerFunc {
	persistenceService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_PERSISTENCE, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	_ = claims
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := getLast30MsgFromDbArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: GetLast30MsgFromDb ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		GetLast30MsgFromDbArgs := persistence.GetLast30MsgFromDbArgs{
			ConversationID: clientArgs.ConversationID,
			Sort:           contant.Sort(clientArgs.Desc),
		}
		GetLast30MsgFromDbReply := persistence.GetLast30MsgFromDbReply{}
		if err := persistenceService.Call(ctx, persistence.SERVICE_GET_LAST30_MSG_FROM_DB, GetLast30MsgFromDbArgs, &GetLast30MsgFromDbReply); err != nil {
			log.Error("Gateway method: GetLast30MsgFromDb Fn persistenceService.Call:persistence.SERVICE_GET_LAST30_MSG_FROM_DB ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(GetLast30MsgFromDbReply))
	}
}

type getLastOneMsgFromDbArgs struct {
	ConversationID string `json:"conversation_id"`
	Desc           bool   `json:"desc"`
}

// 获取会话最后一条消息
func GetLastOneMsgFromDb(ctx context.Context) http.HandlerFunc {
	persistenceService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_PERSISTENCE, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	_ = claims
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := getLastOneMsgFromDbArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: GetLastOneMsgFromDb ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		GetLastOneMsgFromDbArgs := persistence.GetLastOneMsgFromDbArgs{
			ConversationID: clientArgs.ConversationID,
			Sort:           contant.Sort(clientArgs.Desc),
		}
		GetLastOneMsgFromDbReply := persistence.GetLastOneMsgFromDbReply{}
		if err := persistenceService.Call(ctx, persistence.SERVICE_GET_LASTONE_MSG, GetLastOneMsgFromDbArgs, &GetLastOneMsgFromDbReply); err != nil {
			log.Error("Gateway method: GetLastOneMsgFromDb Fn persistenceService.Call:persistence.SERVICE_GET_LASTONE_MSG ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(GetLastOneMsgFromDbReply))
	}
}

type getThe30MsgAfterTheIdArgs struct {
	ConversationID string `json:"conversation_id"`
	MsgID          string `json:"msg_id"`
	Desc           bool   `json:"desc"`
}

// 获取指定会话某个消息id之后的30条消息
func GetThe30MsgAfterTheId(ctx context.Context) http.HandlerFunc {
	persistenceService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_PERSISTENCE, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	_ = claims
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := getThe30MsgAfterTheIdArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: GetThe30MsgAfterTheId ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		GetThe30MsgAfterTheIdArgs := persistence.GetThe30MsgAfterTheIdArgs{
			ConversationID: clientArgs.ConversationID,
			MsgId:          clientArgs.MsgID,
			Sort:           contant.Sort(clientArgs.Desc),
		}
		GetThe30MsgAfterTheIdReply := persistence.GetThe30MsgAfterTheIdReply{}
		if err := persistenceService.Call(ctx, persistence.SERVICE_GET_THE_30MSG_AFTER_THE_ID, GetThe30MsgAfterTheIdArgs, &GetThe30MsgAfterTheIdReply); err != nil {
			log.Error("Gateway method: GetThe30MsgAfterTheId Fn persistenceService.Call:persistence.SERVICE_GET_THE_30MSG_AFTER_THE_ID ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(GetThe30MsgAfterTheIdReply))
	}
}

type getThe30MsgBeforeTheIdArgs struct {
	ConversationID string `json:"conversation_id"`
	MsgID          string `json:"msg_id"`
	Desc           bool   `json:"desc"`
}

// 获取指定会话某个消息id之前的30条消息
func GetThe30MsgBeforeTheId(ctx context.Context) http.HandlerFunc {
	persistenceService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_PERSISTENCE, &rpcx.RpcxClientWithOpt{})
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	_ = claims
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := getThe30MsgBeforeTheIdArgs{}
		encoder := json.NewEncoder(w)
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Error("Gateway method: GetThe30MsgBeforeTheId ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrHttpInvaildParam)
			return
		}
		GetThe30MsgBeforeTheIdArgs := persistence.GetThe30MsgBeforeTheIdArgs{
			ConversationID: clientArgs.ConversationID,
			MsgId:          clientArgs.MsgID,
			Sort:           contant.Sort(clientArgs.Desc),
		}
		GetThe30MsgBeforeTheIdReply := persistence.GetThe30MsgBeforeTheIdReply{}
		if err := persistenceService.Call(ctx, persistence.SERVICE_GET_THE_30MSG_BEFORE_THE_ID, GetThe30MsgBeforeTheIdArgs, &GetThe30MsgBeforeTheIdReply); err != nil {
			log.Error("Gateway method: GetThe30MsgBeforeTheId Fn persistenceService.Call:persistence.SERVICE_GET_THE_30MSG_BEFORE_THE_ID ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		_ = encoder.Encode(quickerr.HttpResponeWarp(GetThe30MsgBeforeTheIdReply))
	}
}
