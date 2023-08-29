package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/quick-im/quick-im-core/services/msghub"
	ser "github.com/quick-im/quick-im-core/services/msghub"
	"github.com/smallnest/rpcx/client"
)

func TestSendMsg(t *testing.T) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8019", "")
	if err != nil {
		t.Error(err)
	}
	xclient := client.NewXClient(ser.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	msg := msghub.SendMsgArgs{
		MsgId:          uuid.New().String(),
		FromSession:    "0",
		ConversationID: "87ba7679-b682-47e7-8499-0385dda22b66",
		MsgType:        0,
		Content:        []byte("哈哈哈哈哈1111"),
		SendTime:       time.Now(),
	}
	reply := msghub.SendMsgReply{}
	err = xclient.Call(context.Background(), ser.SERVICE_SEND_MSG, msg, &reply)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkSendMsg(b *testing.B) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8019", "")
	if err != nil {
		b.Error(err)
	}
	xclient := client.NewXClient(ser.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	msg := msghub.SendMsgArgs{
		MsgId:          uuid.New().String(),
		FromSession:    "0",
		ConversationID: "87ba7679-b682-47e7-8499-0385dda22b66",
		MsgType:        0,
		Content:        []byte("哈哈哈哈哈1111"),
		SendTime:       time.Now(),
	}
	reply := msghub.SendMsgReply{}
	for i := 0; i < b.N; i++ {
		msg.MsgId = uuid.New().String()
		err = xclient.Call(context.Background(), ser.SERVICE_SEND_MSG, msg, &reply)
		if err != nil {
			b.Error(err)
		}
	}
}
