package poll

import (
	"context"
	"net/http"
)

type pollProtoc struct {
	clients map[string]map[uint8]struct{}
}

func InitProtoc() *pollProtoc {
	return &pollProtoc{
		clients: make(map[string]map[uint8]struct{}),
	}
}

func (p *pollProtoc) Handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented") // TODO: Implement
	}
}
