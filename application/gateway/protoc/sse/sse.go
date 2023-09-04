package sse

import (
	"context"
	"net/http"

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

func (s *sseProtoc) Handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			println("The client is disconnected here")
		}()
	}
}
