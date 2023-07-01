package main

import (
	"github.com/quick-im/quick-im-core/services/persistence"
)

func main() {
	if err := persistence.NewServer("0.0.0.0", 8019).Start(); err != nil {
		panic(err)
	}
}
