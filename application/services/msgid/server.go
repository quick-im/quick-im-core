package main

import "github.com/quick-im/quick-im-core/services/msgid"

func main() {
	if err := msgid.NewServer("0.0.0.0", 8018).Start(); err != nil {
		panic(err)
	}
}
