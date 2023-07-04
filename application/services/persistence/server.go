package main

import (
	"github.com/quick-im/quick-im-core/services/persistence"
)

func main() {
	if err := persistence.NewServer(
		persistence.WithIp("0.0.0.0"),
		persistence.WithPort(8015),
	).Start(); err != nil {
		panic(err)
	}
}
