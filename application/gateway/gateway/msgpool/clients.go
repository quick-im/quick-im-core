package msgpool

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/quick-im/quick-im-core/internal/codec"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
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

var Gid = uuid.New().String()

var cs = &clients{
	gid: Gid,
	ch:  make(chan *clientAndCh),
}

type chWarp struct {
	sid      string
	platform uint8
	ch       chan model.Msg
}

var subs = make(map[string]map[uint8]chWarp)
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
	lock.Lock()
	if _, ok := subs[sid]; !ok {
		subs[sid] = make(map[uint8]chWarp)
	}
	subs[sid][platform] = chWarp{
		sid:      sid,
		platform: platform,
		ch:       make(chan model.Msg),
	}
	lock.Unlock()
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

func GetMsgChannel(sid string, platform uint8) (ch chWarp, ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	if chs, ok := subs[sid]; ok {
		if ch, ok := chs[platform]; ok {
			return ch, ok
		}
	}
	return chWarp{}, false
}

func (cch chWarp) UnRegistry() {
	lock.Lock()
	defer lock.Unlock()
	if chs, ok := subs[cch.sid]; ok {
		if c, ok := chs[cch.platform]; ok {
			close(c.ch)
			delete(chs, cch.platform)
		}
		if len(chs) == 0 {
			delete(subs, cch.sid)
			println("unregister", cch.sid, "-", cch.platform)
		}
	}

}

func (cch chWarp) GetCh() <-chan model.Msg {
	return cch.ch
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
		lock.RLock()
		// 消息分发
		for i := range msgData.ToSessions {
			if ch, ok := subs[msgData.ToSessions[i].SessionId][msgData.ToSessions[i].Platform]; ok {
				ch.ch <- msgData.MetaData
			}
		}
		lock.RUnlock()
	}
}
