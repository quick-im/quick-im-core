package main

import (
	"github.com/quick-im/quick-im-core/services/msghub"
)

func main() {
	if err := msghub.NewServer("0.0.0.0", 8019).Start(); err != nil {
		panic(err)
	}
}
