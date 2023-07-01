package main

import (
	"github.com/quick-im/quick-im-core/services/conversation"
)

func main() {
	if err := conversation.NewServer(
		conversation.SetOptIp("0.0.0.0"),
		conversation.SetOptPort(8016),
	).Start(); err != nil {
		panic(err)
	}
}
