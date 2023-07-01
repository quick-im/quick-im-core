package main

import (
	"github.com/quick-im/quick-im-core/services/msgid"
)

func main() {
	if err := msgid.NewServer(
		msgid.SetOptIp("0.0.0.0"),
		msgid.SetOptPort(8018),
	).Start(); err != nil {
		panic(err)
	}
}
