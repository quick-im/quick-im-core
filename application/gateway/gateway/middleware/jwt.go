package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/jwt"
)

type AuthHandlerFunc func(ctx context.Context) http.HandlerFunc

func JwtAuth(ctx context.Context, h AuthHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		urlToken := r.URL.Query().Get("token")
		if authorization == "" {
			authorization = urlToken
		}
		claims, err := jwt.ParseToken(strings.TrimPrefix(authorization, "Bearer "))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(ctx, contant.HTTP_CTX_JWT_CLAIMS, claims)
		h(ctx)(w, r)
	}
}
