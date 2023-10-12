package test

import (
	"context"
	"testing"

	"github.com/quick-im/quick-im-core/services/conversation"
	"github.com/smallnest/rpcx/client"
)

func TestServer(t *testing.T) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8016", "")
	if err != nil {
		t.Error(err)
	}
	xclient := client.NewXClient(conversation.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := conversation.CreateConversationArgs{
		ConversationType: 0,
		SessionList:      []string{"123456"},
	}
	reply := &conversation.CreateConversationReply{}
	if err := xclient.Call(context.Background(), conversation.SERVICE_CREATE_CONVERSATION, args, reply); err != nil {
		t.Error(err)
	}
	t.Log(reply.ConversationID)
}
