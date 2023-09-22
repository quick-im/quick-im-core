package sse

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

type sseProtoc struct {
	clients map[string]map[uint8]<-chan model.Msg
}

func InitProtoc() *sseProtoc {
	return &sseProtoc{
		clients: make(map[string]map[uint8]<-chan model.Msg),
	}
}

type serverSentEvent struct {
	EventType string
	Data      string
}

func (s *sseProtoc) Handler(ctx context.Context) http.HandlerFunc {
	claims := helper.GetCtxValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, contant.JWTClaimsCtxType)
	var log logger.Logger
	log = helper.GetCtxValue(ctx, contant.CTX_LOGGER_KEY, log)
	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			println("The client is disconnected here")
		}()
		// Send Browser Connection
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		eventData := &serverSentEvent{
			EventType: "message",
			Data:      "",
		}
		ch, ok := msgpool.GetMsgChannel(claims.Sid, claims.Platform)
		if !ok {
			log.Error("PollHandler: msg channel not found")
			return
		}
		for {
			msg := <-ch
			data, err := json.Marshal(msg)
			if err != nil {
				log.Error("PollHandler: marshal msg error", err.Error())
				continue
			}
			fmt.Fprintf(w, "id: %d\n", time.Now().Unix())
			fmt.Fprintf(w, "event: %s\n", eventData.EventType)
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}
}
