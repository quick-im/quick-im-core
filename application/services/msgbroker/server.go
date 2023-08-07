package main

import (
	"context"

	"github.com/quick-im/quick-im-core/services/msgbroker"
)

func main() {
	ctx := context.Background()
	if err := msgbroker.NewServer(
		msgbroker.WithIp("0.0.0.0"),
		msgbroker.WithPort(8017),
	).Start(ctx); err != nil {
		panic(err)
	}
}
