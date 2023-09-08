package msgpool

import (
	"context"
	"testing"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/services/msgbroker"
	cclient "github.com/rpcxio/rpcx-consul/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

func TestRegisterTerm(t *testing.T) {
	go RunMsgPollServer(context.Background())
	ch := make(chan *protocol.Message)
	d, _ := cclient.NewConsulDiscovery(config.ServerPrefix, msgbroker.SERVER_NAME, []string{"127.0.0.1:8500"}, nil)
	xclient := client.NewBidirectionalXClient(msgbroker.SERVER_NAME, client.Failtry, client.ConsistentHash, d, client.DefaultOption, ch)
	defer xclient.Close()
	type args struct {
		ctx      context.Context
		c        client.XClient
		ch       chan *protocol.Message
		sid      string
		platform uint8
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "reg",
			args: args{
				ctx:      context.Background(),
				c:        xclient,
				ch:       ch,
				sid:      "50864896-8136-4a43-8a48-1d3325a7f78f",
				platform: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterTerm(tt.args.ctx, tt.args.c, tt.args.ch, tt.args.sid, tt.args.platform); (err != nil) != tt.wantErr {
				t.Errorf("RegisterTerm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	select {}
}
