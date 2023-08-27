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
		ConversationID: "0.0.0.0",
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
