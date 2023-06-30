package main

import (
	"github.com/quick-im/quick-im-core/services/conversation"
)

func main() {
	if err := conversation.NewServer("0.0.0.0", 8019).Start(); err != nil {
		panic(err)
	}
}
