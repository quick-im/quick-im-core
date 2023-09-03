package sse

import (
	"context"
	"net/http"
)

type sseProtoc struct {
	clients map[string]map[uint8]struct{}
}

func InitProtoc() *sseProtoc {
	return &sseProtoc{
		clients: make(map[string]map[uint8]struct{}),
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
