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
	cs  []clientAndCh
}

type clientAndCh struct {
	c  client.XClient
	ch <-chan *protocol.Message
}

var cs = &clients{
	gid: uuid.New().String(),
	cs:  make([]clientAndCh, 0),
}

func RegisterTerm(c client.XClient, ch <-chan *protocol.Message, sid string, platform uint8) error {
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
		cs.cs = append(cs.cs, clientAndCh{
			c:  c,
			ch: ch,
		})
	} else {
		_ = c.Close()
	}
	return nil
}
