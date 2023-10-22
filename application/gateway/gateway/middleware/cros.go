package middleware

import (
	"context"
	"net/http"
)

func AllowCros(ctx context.Context, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Content-Type", "application/json")
		h(w, r)
	}
}
