package main

import (
	"context"

	"github.com/quick-im/quick-im-core/services/msghub"
)

func main() {
	ctx := context.Background()
	if err := msghub.NewServer(
		msghub.WithIp("0.0.0.0"),
		msghub.WithPort(8019),
		msghub.WithConsulServers("nats://127.0.0.1:4222"),
	).Start(ctx); err != nil {
		panic(err)
	}
}
