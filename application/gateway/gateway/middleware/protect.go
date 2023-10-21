package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/quickerr"
)

type ProtectHandlerFunc func(ctx context.Context) http.HandlerFunc

// 用于通过IP白名单保护一些接口
func ProtectApi(ctx context.Context, h ProtectHandlerFunc) http.HandlerFunc {
	ipwhite := helper.GetCtxValue(ctx, contant.CTX_IP_WHITELIST_KEY, contant.IPWhiteListCtxType{})
	CIDRList := make([]*net.IPNet, 0, len(ipwhite))
	for i := range ipwhite {
		var cidr string
		cidr = ipwhite[i]
		sp := strings.IndexByte(cidr, '/')
		if sp < 0 {
			if strings.IndexByte(cidr, '.') < 0 {
				cidr = fmt.Sprintf("%s/24", cidr)
			} else {
				cidr = fmt.Sprintf("%s/128", cidr)
			}
			_, c, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			CIDRList[i] = c
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		coder := json.NewEncoder(w)
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var pass bool
		for i := range ipwhite {
			if CIDRList[i] != nil {
				CIDRList[i].Contains(net.IP(host))
			}
			pass = true
			break
		}
		if !pass {
			coder.Encode(quickerr.ErrNotAllowedRequest)
		}
		h(ctx)(w, r)
	}
}
