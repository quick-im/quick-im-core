package protoc

import (
	"context"
	"net/http"
)

type ProtocHandler interface {
	Handler(ctx context.Context) http.HandlerFunc
}
