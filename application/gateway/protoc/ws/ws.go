package ws

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/quick-im/quick-im-core/application/gateway/gateway/msgpool"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsProtoc struct {
	clients map[string]map[uint8]<-chan model.Msg
}

func InitProtoc() *wsProtoc {
	return &wsProtoc{
		clients: make(map[string]map[uint8]<-chan model.Msg),
	}
}

func (ws *wsProtoc) Handler(ctx context.Context) http.HandlerFunc {
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Gateway method: WsHandler upgrader.Upgrade,err: ", err.Error())
			return
		}
		defer conn.Close()
		chWarp, ok := msgpool.GetMsgChannel(claims.Sid, claims.Platform)
		if !ok {
			log.Error("WsHandler: msg channel not found")
			return
		}
		ch := chWarp.GetCh()
		for {
			msg := ch
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Error("Gateway method: WsHandler conn.WriteJSON,err: ", err.Error())
				continue
			}
		}
	}
}
