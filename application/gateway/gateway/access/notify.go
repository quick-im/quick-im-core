package access

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/quick-im/quick-im-core/application/gateway/protoc"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/logger"
)

func NotifyHandler(ctx context.Context) http.HandlerFunc {
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	_ = claims
	_ = log
	return func(w http.ResponseWriter, r *http.Request) {
		encode := json.NewEncoder(w)
		clientProtoc := r.URL.Query().Get("protoc")
		handler, err := protoc.Handler(clientProtoc)
		if err != nil {
			_ = encode.Encode(err)
			return
		}
		// 将按照客户端指定的协议处理
		handler.Handler(ctx)(w, r)
	}
}
