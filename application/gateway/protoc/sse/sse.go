package sse

import (
	"context"
	"net/http"

	"github.com/r3labs/sse/v2"
)

var server = sse.New()

type SSEProtoc struct {
}

func (s *SSEProtoc) Handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			println("The client is disconnected here")
		}()

		server.ServeHTTP(w, r)
	}
}
