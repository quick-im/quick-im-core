package poll

import (
	"context"
	"net/http"

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
	return func(w http.ResponseWriter, r *http.Request) {
		panic("not implemented") // TODO: Implement
	}
}
