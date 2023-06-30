package main

import "github.com/quick-im/quick-im-core/services/mid"

func main() {
	if err := mid.NewServer("0.0.0.0", 8018).Start(); err != nil {
		panic(err)
	}
}
