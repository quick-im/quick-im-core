package messaging

import (
	"log"
	"strings"

	"github.com/nats-io/nats.go"
)

type natsClientOpt struct {
	servers     []string
	useJsStream bool
}

type NatsWarp struct {
	nc          *nats.Conn
	useJsStream bool
	js          nats.JetStreamContext
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

func WithJetStream(jetstream bool) natsOpt {
	return func(nco *natsClientOpt) {
		nco.useJsStream = jetstream
	}
}

func (n *natsClientOpt) GetNats() *NatsWarp {
	nc, err := nats.Connect(strings.Join(n.servers, ","))
	if err != nil {
		log.Fatal(err)
	}
	return &NatsWarp{
		nc,
		n.useJsStream,
		nil,
	}
}

func (n *NatsWarp) JetStream(opts ...nats.JSOpt) (nats.JetStreamContext, error) {
	js, err := n.nc.JetStream(opts...)
	if err != nil {
		return nil, err
	}
	n.js = js
	n.useJsStream = true
	return js, err
}

func (n *NatsWarp) Close() {
	n.nc.Close()
}
