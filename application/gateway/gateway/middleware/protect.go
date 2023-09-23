package middleware

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/quickerr"
)

type ProtectHandlerFunc func(ctx context.Context) http.HandlerFunc

// 用于通过IP白名单保护一些接口
func ProtectApi(ctx context.Context, h ProtectHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coder := json.NewEncoder(w)
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var pass bool
		for i := range config.IPWhite {
			config.IPWhite[i].Contains(net.IP(host))
			pass = true
			break
		}
		if !pass {
			coder.Encode(quickerr.ErrNotAllowedRequest)
		}
		h(ctx)(w, r)
	}
}
