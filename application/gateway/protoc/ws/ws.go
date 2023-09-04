package ws

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		_ = conn
	}
}
