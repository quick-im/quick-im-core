package main

import (
	"github.com/quick-im/quick-im-core/services/persistence"
)

func main() {
	if err := persistence.NewServer(
		persistence.SetOptIp("0.0.0.0"),
		persistence.SetOptPort(8015),
	).Start(); err != nil {
		panic(err)
	}
}
