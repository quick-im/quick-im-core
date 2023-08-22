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
		msgbroker.WithNatsServers("nats://127.0.0.1:4222"),
		msgbroker.WithUseConsulRegistry(true),
		msgbroker.WithConsulServers("127.0.0.1:8500"),
	).Start(ctx); err != nil {
		panic(err)
	}
}
