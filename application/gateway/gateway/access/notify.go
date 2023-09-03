package access

import (
	"context"
	"net/http"
)

func NotifyHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
