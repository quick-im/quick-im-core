package messaging

import (
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

type natsClientOpt struct {
	servers []string
}

type natsOpt func(*natsClientOpt)

func NewNatsWithOpt(opts ...natsOpt) *natsClientOpt {
	n := &natsClientOpt{
		servers: make([]string, 0),
	}
	for i := range opts {
		opts[i](n)
	}
	return n
}

func WithServer(server string) natsOpt {
	return func(nco *natsClientOpt) {
		nco.servers = append(nco.servers, server)
	}
}

func WithServers(servers ...string) natsOpt {
	return func(nco *natsClientOpt) {
		nco.servers = servers
	}
}

func (n *natsClientOpt) GetNats() *nats.Conn {
	nc, err := nats.Connect(strings.Join(n.servers, ","))
	if err != nil {
		log.Fatal(err)
	}
	return nc
}
