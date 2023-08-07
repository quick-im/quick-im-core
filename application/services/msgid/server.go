package main

import (
	"context"

	"github.com/quick-im/quick-im-core/services/msgid"
)

func main() {
	ctx := context.Background()
	if err := msgid.NewServer(
		msgid.WithIp("0.0.0.0"),
		msgid.WithPort(8018),
		msgid.WithOpenTracing(true),
		msgid.WithJeagerServiceName("msgid"),
		msgid.WithJeagerAgentHostPort("127.0.0.1:6831"),
	).Start(ctx); err != nil {
		panic(err)
	}
}
