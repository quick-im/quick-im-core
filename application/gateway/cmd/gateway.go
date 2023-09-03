package main

import (
	"context"
	"log"

	"github.com/quick-im/quick-im-core/application/gateway/api/server"
)

func main() {
	log.SetFlags(log.Llongfile)
	ctx := context.Background()
	ser := server.NewApiServer(
		server.WithIp("0.0.0.0"),
		server.WithPort(8088),
	)
	ser.InitAndStartServer(ctx)
}
