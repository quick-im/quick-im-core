package main

import (
	"context"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/services/msghub"
)

func main() {
	ctx := context.Background()
	if err := msghub.NewServer(
		config.WithIp("0.0.0.0"),
		config.WithPort(8019),
		config.WithOpenTracing(true),
		config.WithJeagerAgentHostPort("127.0.0.1:6831"),
		config.WithNatsServers("nats://127.0.0.1:4222"),
		config.WithUseConsulRegistry(true),
		config.WithConsulServers("127.0.0.1:8500"),
	).Start(ctx); err != nil {
		panic(err)
	}
}
