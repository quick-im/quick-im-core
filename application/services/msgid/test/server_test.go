package main

import (
	"context"
	"testing"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/quick-im/quick-im-core/internal/tracing/plugin"
	"github.com/quick-im/quick-im-core/services/msgid"
	cclient "github.com/rpcxio/rpcx-consul/client"
	"github.com/smallnest/rpcx/client"
)

func TestServer(t *testing.T) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8018", "")
	if err != nil {
		t.Error(err)
	}
	xclient := client.NewXClient(msgid.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	plugins := client.NewPluginContainer()
	tracer, ctx := plugin.AddClientTrace("client", "127.0.0.1:6831", plugins)
	defer tracer.Shutdown(ctx)
	xclient.SetPlugins(plugins)
	args := msgid.GenerateMessageIDArgs{
		ConversationID:   "123",
		ConversationType: 3,
	}
	reply := &msgid.GenerateMessageIDReply{}
	if err := xclient.Call(ctx, msgid.SERVICE_GENERATE_MESSAGE_ID, args, reply); err != nil {
		t.Error(err)
	}
	t.Log(reply.MsgID)
}

func TestConsul(t *testing.T) {
	d, _ := cclient.NewConsulDiscovery(config.ServerPrefix, msgid.SERVER_NAME, []string{"127.0.0.1:8500"}, nil)
	xclient := client.NewXClient(msgid.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	plugins := client.NewPluginContainer()
	tracer, ctx := plugin.AddClientTrace("client", "127.0.0.1:6831", plugins)
	defer tracer.Shutdown(ctx)
	xclient.SetPlugins(plugins)
	args := msgid.GenerateMessageIDArgs{
		ConversationID:   "123",
		ConversationType: 3,
	}
	reply := &msgid.GenerateMessageIDReply{}
	if err := xclient.Call(ctx, msgid.SERVICE_GENERATE_MESSAGE_ID, args, reply); err != nil {
		t.Error(err)
	}
	t.Log(reply.MsgID)
}

func TestOptClient(t *testing.T) {
	c, err := rpcx.NewClient(
		rpcx.WithBasePath(config.ServerPrefix),
		rpcx.WithServerAddress("127.0.0.1:8018"),
		rpcx.WithServiceName(msgid.SERVER_NAME),
		rpcx.WithOpenTracing(true),
		rpcx.WithJeagerAgentHostPort("127.0.0.1:6831"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer c.CloseAndShutdownTrace()
	args := msgid.GenerateMessageIDArgs{
		ConversationID:   "123",
		ConversationType: 3,
	}
	reply := &msgid.GenerateMessageIDReply{}
	if err := c.Call(context.Background(), msgid.SERVICE_GENERATE_MESSAGE_ID, args, reply); err != nil {
		t.Error(err)
	}
	t.Log(reply.MsgID)
}

func BenchmarkServer(b *testing.B) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8018", "")
	if err != nil {
		b.Error(err)
	}
	xclient := client.NewXClient(msgid.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := msgid.GenerateMessageIDArgs{
		ConversationID:   "123",
		ConversationType: 3,
	}
	reply := &msgid.GenerateMessageIDReply{}
	for i := 0; i < b.N; i++ {
		if err := xclient.Call(context.Background(), msgid.SERVICE_GENERATE_MESSAGE_ID, args, reply); err != nil {
			b.Error(err)
		}
	}
}
