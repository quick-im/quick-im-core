package poll

import (
	"context"
	"net/http"
)

type PollProtoc struct {
}

func (p *PollProtoc) Handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented") // TODO: Implement
	}
}
