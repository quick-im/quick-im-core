package test

import (
	"context"
	"testing"
	"time"

	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"github.com/quick-im/quick-im-core/services/persistence"
	"github.com/smallnest/rpcx/client"
)

func TestServerSaveMsgToDb(t *testing.T) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8015", "")
	if err != nil {
		t.Error(err)
	}
	xclient := client.NewXClient(persistence.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	args := persistence.SaveMsgToDbArgs{
		Msgs: []model.Msg{
			{
				MsgId:          "EBVYE-795J-246S-RBG",
				ConversationID: "1",
				FromSession:    "0",
				SendTime:       time.Now(),
				Status:         0,
				Type:           0,
				Content:        "hihi哈哈",
			},
			{
				MsgId:          "EBVYE-796J-266S-RBG",
				ConversationID: "1",
				FromSession:    "0",
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

func TestServerGetMsg(t *testing.T) {
	d, err := client.NewPeer2PeerDiscovery("tcp@127.0.0.1:8015", "")
	if err != nil {
		t.Error(err)
	}
	xclient := client.NewXClient(persistence.SERVER_NAME, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	t.Run("SERVICE_GET_MSG_FROM_DB_IN_RANGE", func(t *testing.T) {
		reply := &persistence.GetMsgFromDbInRangeReply{}
		err := xclient.Call(context.Background(), persistence.SERVICE_GET_MSG_FROM_DB_IN_RANGE, persistence.GetMsgFromDbInRangeArgs{
			ConversationID: "1",
			StartMsgId:     "EBVYE-795J-246S-RBG",
			EndMsgId:       "EBVYE-796J-266S-RBG",
			Sort:           false,
		}, reply)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", reply.Msg)
	})
	t.Run("SERVICE_GET_LAST30_MSG_FROM_DB", func(t *testing.T) {
		reply := &persistence.GetLast30MsgFromDbReply{}
		err := xclient.Call(context.Background(), persistence.SERVICE_GET_LAST30_MSG_FROM_DB, persistence.GetLast30MsgFromDbArgs{
			ConversationID: "1",
			Sort:           false,
		}, reply)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", reply.Msg)
	})
	t.Run("SERVICE_GET_THE_30MSG_BEFORE_THE_ID", func(t *testing.T) {
		reply := &persistence.GetThe30MsgBeforeTheIdReply{}
		err := xclient.Call(context.Background(), persistence.SERVICE_GET_THE_30MSG_BEFORE_THE_ID, persistence.GetThe30MsgBeforeTheIdArgs{
			ConversationID: "1",
			MsgId:          "EBVYE-796J-266S-RBG",
			Sort:           false,
		}, reply)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", reply.Msg)
	})
	t.Run("SERVICE_GET_THE_30MSG_AFTER_THE_ID", func(t *testing.T) {
		reply := &persistence.GetThe30MsgAfterTheIdReply{}
		err := xclient.Call(context.Background(), persistence.SERVICE_GET_THE_30MSG_AFTER_THE_ID, persistence.GetThe30MsgAfterTheIdArgs{
			ConversationID: "1",
			MsgId:          "EBVYE-795J-266S-RBG",
			Sort:           false,
		}, reply)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", reply.Msg)
	})

}
