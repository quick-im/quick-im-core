package poll

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/quick-im/quick-im-core/application/gateway/gateway/msgpool"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

type pollProtoc struct {
	clients map[string]map[uint8]<-chan model.Msg
}

func InitProtoc() *pollProtoc {
	return &pollProtoc{
		clients: make(map[string]map[uint8]<-chan model.Msg),
	}
}

func (p *pollProtoc) Handler(ctx context.Context) http.HandlerFunc {
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		chWarp, ok := msgpool.GetMsgChannel(claims.Sid, claims.Platform)
		if !ok {
			log.Error("PollHandler: msg channel not found")
			return
		}
		defer chWarp.UnRegistry()
		ch := chWarp.GetCh()
		timer := time.NewTimer(time.Second * 30)
		defer timer.Stop()
		select {
		case <-timer.C:
			w.WriteHeader(http.StatusNoContent)
		case <-r.Context().Done():
			w.WriteHeader(http.StatusNoContent)
		case msg := <-ch:
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fmt.Println(msg)
			_ = encoder.Encode(msg)
			return
		}
	}
}
