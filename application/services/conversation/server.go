package main

import (
	"github.com/quick-im/quick-im-core/services/conversation"
)

func main() {
	if err := conversation.NewServer(
		conversation.WithIp("0.0.0.0"),
		conversation.WithPort(8016),
	).Start(); err != nil {
		panic(err)
	}
}
