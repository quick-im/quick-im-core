package main

import (
	"github.com/quick-im/quick-im-core/services/msgbroker"
)

func main() {
	if err := msgbroker.NewServer("0.0.0.0", 8019).Start(); err != nil {
		panic(err)
	}
}
