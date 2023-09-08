package access

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/quickerr"
	"github.com/quick-im/quick-im-core/internal/rpcx"
)

type sendMsgArgs struct {
	ConversationID string `json:"conversation_id"`
	Status         int32  `json:"status"`
	Type           int32  `json:"type"`
	Content        string `json:"content"`
}

func SendMsgHandler(ctx context.Context) http.HandlerFunc {
	msghubService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGHUB, &rpcx.RpcxClientWithOpt{})
	msgidService := helper.GetCtxValue(ctx, contant.CTX_SERVICE_MSGID, &rpcx.RpcxClientWithOpt{})
	_ = msghubService
	_ = msgidService
	return func(w http.ResponseWriter, r *http.Request) {
		clientArgs := sendMsgArgs{}
		if err := json.NewDecoder(r.Body).Decode(&clientArgs); err != nil {
			log.Println(err)
			_ = json.NewEncoder(w).Encode(quickerr.ErrHttpInvaildParam)
			return
		}
	}
}
