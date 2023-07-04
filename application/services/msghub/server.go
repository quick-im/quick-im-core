package main

import (
	"github.com/quick-im/quick-im-core/services/msghub"
)

func main() {
	if err := msghub.NewServer(
		msghub.WithIp("0.0.0.0"),
		msghub.WithPort(8019),
	).Start(); err != nil {
		panic(err)
	}
}
