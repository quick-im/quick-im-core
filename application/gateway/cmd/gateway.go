package main

import (
	"context"
	"net/http"

	"github.com/quick-im/quick-im-core/application/gateway/api"
)

func main() {
	ctx := context.Background()
	http.HandleFunc("/notify", api.NotifyHandler(ctx))
	http.ListenAndServe(":8080", nil)
}
