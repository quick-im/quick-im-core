package messaging

import (
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

func GetNats() *nats.Conn {
	servers := []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"}
	nc, err := nats.Connect(strings.Join(servers, ","))
	if err != nil {
		log.Fatal(err)
	}
	return nc
}
