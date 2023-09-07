package msgpool

import (
	"context"

	"github.com/google/uuid"
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

func RegisterTerm(ctx context.Context, c client.XClient, ch <-chan *protocol.Message, sid string, platform uint8) error {
	regSessionArgs := msgbroker.RegisterSessionInfo{
		Platform:    platform,
		GatewayUuid: cs.gid,
		SessionId:   sid,
	}
	regSessionReply := msgbroker.RegisterSessionReply{}
	if err := c.Call(context.Background(), msgbroker.SERVICE_REGISTER_SESSION, regSessionArgs, &regSessionReply); err != nil {
		return err
	}
	if regSessionReply.NeedKeep {
		// 保持连接
		select {
		case cs.ch <- &clientAndCh{
			c:  c,
			ch: make(chan *protocol.Message),
		}:
		case <-ctx.Done():
			println("register timeout", sid, "-", platform)
		}
	} else {
		_ = c.Close()
	}
	return nil
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
	for msg := range cn.ch {
		_ = msg
	}
}
