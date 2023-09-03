package middleware

import (
	"context"
	"net/http"
)

func JwtAuth(ctx context.Context, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	}
}
