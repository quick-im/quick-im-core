package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/quick-im/quick-im-core/internal/config"
	param "github.com/quick-im/quick-im-core/internal/quickparam/msgbroker"
	"github.com/quick-im/quick-im-core/services/msgbroker"
	cclient "github.com/rpcxio/rpcx-consul/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

func TestRecvMsg(t *testing.T) {
	ch := make(chan *protocol.Message)
	d, _ := cclient.NewConsulDiscovery(config.ServerPrefix, msgbroker.SERVER_NAME, []string{"127.0.0.1:8500"}, nil)
	xclient := client.NewBidirectionalXClient(msgbroker.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption, ch)
	defer xclient.Close()
	gatewayUuid := uuid.New().String()
	args := param.RegisterSessionInfo{
		Platform:    0,
		GatewayUuid: gatewayUuid,
		SessionId:   "1",
	}
	reply := &param.RegisterSessionReply{}

	if err := xclient.Call(context.Background(), msgbroker.SERVICE_REGISTER_SESSION, args, reply); err != nil {
		t.Error(err)
	}
	for msg := range ch {
		fmt.Printf("receive msg from server: %s\n", msg.Payload)
	}
}
