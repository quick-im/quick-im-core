package test

import (
	"context"
	"testing"
	"time"

	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"github.com/quick-im/quick-im-core/services/persistence"
	"github.com/smallnest/rpcx/client"
)

func TestServer(t *testing.T) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8015", "")
	if err != nil {
		t.Error(err)
	}
	xclient := client.NewXClient(persistence.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := persistence.SaveMsgToDbArgs{
		Msgs: []model.Msg{
			{
				MsgId:          "EBVYE-785J-246S-RBG",
				ConvercationID: "1",
				FromSession:    0,
				SendTime:       time.Now(),
				Status:         0,
				Type:           0,
				Content:        "hihi哈哈",
			},
			{
				MsgId:          "EBVYE-785J-266S-RBG",
				ConvercationID: "1",
				FromSession:    0,
				SendTime:       time.Now(),
				Status:         0,
				Type:           0,
				Content:        "hihi哈哈",
			},
		},
	}
	reply := &persistence.SaveMsgToDbReply{}
	if err := xclient.Call(context.Background(), persistence.SERVICE_SAVE_MSG_TO_DB, args, reply); err != nil {
		t.Error(err)
	}
	t.Log(reply)
}
