package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/services/msgbroker"
	cclient "github.com/rpcxio/rpcx-consul/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

func TestRecvMsg(t *testing.T) {
	msgData := msgbroker.BroadcastMsgWarp{}
	c := codec.GobUtils[msgbroker.BroadcastMsgWarp]{}
	ch := make(chan *protocol.Message)
	d, _ := cclient.NewConsulDiscovery(config.ServerPrefix, msgbroker.SERVER_NAME, []string{"127.0.0.1:8500"}, nil)
	xclient := client.NewBidirectionalXClient(msgbroker.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption, ch)
	defer xclient.Close()
	gatewayUuid := uuid.New().String()
	args := msgbroker.RegisterSessionInfo{
		Platform:    0,
		GatewayUuid: gatewayUuid,
		SessionId:   "50864896-8136-4a43-8a48-1d3325a7f78f",
	}
	reply := &msgbroker.RegisterSessionReply{}

	if err := xclient.Call(context.Background(), msgbroker.SERVICE_REGISTER_SESSION, args, reply); err != nil {
		t.Error(err)
	}
	for msg := range ch {
		// fmt.Printf("receive msg from server: %s\n", msg.Payload)
		c.Decode(msg.Payload, &msgData)
		fmt.Printf("receive msg from server and decode: %#v\n", msgData)
	}
}
