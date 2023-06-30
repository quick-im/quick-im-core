package main

import (
	"context"
	"testing"

	"github.com/quick-im/quick-im-core/services/mid"
	"github.com/smallnest/rpcx/client"
)

func TestServer(t *testing.T) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8018", "")
	if err != nil {
		t.Error(err)
	}
	xclient := client.NewXClient(mid.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := mid.GenerateMessageIDArgs{
		ConversationID:   "123",
		ConversationType: 3,
	}
	reply := &mid.GenerateMessageIDReply{}
	if err := xclient.Call(context.Background(), mid.SERVICE_GENERATE_MESSAGE_ID, args, reply); err != nil {
		t.Error(err)
	}
	t.Log(reply.MsgID)
}

func BenchmarkServer(b *testing.B) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8018", "")
	if err != nil {
		b.Error(err)
	}
	xclient := client.NewXClient(mid.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := mid.GenerateMessageIDArgs{
		ConversationID:   "123",
		ConversationType: 3,
	}
	reply := &mid.GenerateMessageIDReply{}
	for i := 0; i < b.N; i++ {
		if err := xclient.Call(context.Background(), mid.SERVICE_GENERATE_MESSAGE_ID, args, reply); err != nil {
			b.Error(err)
		}
	}
}
