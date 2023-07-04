package main

import (
	"github.com/quick-im/quick-im-core/services/msgbroker"
)

func main() {
	if err := msgbroker.NewServer(
		msgbroker.WithIp("0.0.0.0"),
		msgbroker.WithPort(8017),
	).Start(); err != nil {
		panic(err)
	}
}
