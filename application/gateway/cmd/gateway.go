package main

import (
	"context"
	"log"

	"github.com/quick-im/quick-im-core/application/gateway/gateway/msgpool"
	"github.com/quick-im/quick-im-core/application/gateway/gateway/server"
)

func main() {
	log.SetFlags(log.Llongfile)
	ctx := context.Background()
	ser := server.NewApiServer(
		server.WithIp("0.0.0.0"),
		server.WithPort(8088),
	)
	go msgpool.RunMsgPollServer(ctx)
	ser.InitAndStartServer(ctx)
}
