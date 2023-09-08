package msgpool

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/services/msgbroker"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

type clients struct {
	gid string
	ch  chan *clientAndCh
}

type clientAndCh struct {
	c  client.XClient
	ch chan *protocol.Message
}

var cs = &clients{
	gid: uuid.New().String(),
	ch:  make(chan *clientAndCh),
}

var subs = make(map[string]map[uint8]struct{})
var lock = sync.RWMutex{}

// 如果needKeep为true，请不要在外部关闭XClient
func RegisterTerm(ctx context.Context, c client.XClient, ch chan *protocol.Message, sid string, platform uint8) (needKeep bool, err error) {
	regSessionArgs := msgbroker.RegisterSessionInfo{
		Platform:    platform,
		GatewayUuid: cs.gid,
		SessionId:   sid,
	}
	regSessionReply := msgbroker.RegisterSessionReply{}
	if err := c.Call(context.Background(), msgbroker.SERVICE_REGISTER_SESSION, regSessionArgs, &regSessionReply); err != nil {
		return needKeep, err
	}
	// println(regSessionReply.NeedKeep)
	if regSessionReply.NeedKeep {
		// 保持连接
		select {
		case cs.ch <- &clientAndCh{
			c:  c,
			ch: ch,
		}:
			needKeep = true
		case <-ctx.Done():
			needKeep = false
			println("register timeout", sid, "-", platform)
		}
	}
	return needKeep, nil
}

func RunMsgPollServer(ctx context.Context) {
	cs.Run(ctx)
}

func (c *clients) Run(ctx context.Context) {
	for {
		if cn, ok := <-c.ch; ok {
			go cn.ListenMsg(ctx)
		}
	}
}

func (cn *clientAndCh) ListenMsg(ctx context.Context) {
	defer cn.c.Close()
	msgData := msgbroker.BroadcastMsgWarp{}
	codec := codec.GobUtils[msgbroker.BroadcastMsgWarp]{}
	for msg := range cn.ch {
		if err := codec.Decode(msg.Payload, &msgData); err != nil {
			println("Decode Msg Failed: ", err)
			continue
		}
	}
}
