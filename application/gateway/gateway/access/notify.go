package access

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/quick-im/quick-im-core/application/gateway/gateway/msgpool"
	"github.com/quick-im/quick-im-core/application/gateway/protoc"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/quickerr"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/smallnest/rpcx/protocol"
)

func NotifyHandler(ctx context.Context) http.HandlerFunc {
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	var log logger.Logger
	_ = claims
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	msgbrokerService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGBORKER, &rpcx.RpcxClientWithOpt{})
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		clientProtoc := r.URL.Query().Get("protoc")
		handler, err := protoc.Handler(clientProtoc)
		if err != nil {
			_ = encoder.Encode(err)
			return
		}
		msgbrokerServiceClient, err := msgbrokerService.GetOnce()
		if err != nil {
			log.Error("Gateway method: NotifyHandler msgbrokerService.GetOnce() ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		ch := make(chan *protocol.Message)
		keep, err := msgpool.RegisterTerm(ctx, msgbrokerServiceClient, ch, claims.Sid, claims.Platform)
		if err != nil {
			msgbrokerServiceClient.Close()
			close(ch)
			log.Error("Gateway method: NotifyHandler msgpool.RegisterTerm() ,err: ", err.Error())
			_ = encoder.Encode(quickerr.ErrInternalServiceCallFailed)
			return
		}
		// 每个gateway只需通过一个channel消息池，所以不需要重复监听
		if !keep {
			msgbrokerServiceClient.Close()
			close(ch)
		}
		// 将按照客户端指定的协议处理
		handler.Handler(ctx)(w, r)
	}
}
