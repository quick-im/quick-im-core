package msgpool

import (
	"context"
	"sync"
	"time"

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
	queue    *ConcurrentQueue[model.Msg]
}

var subs = make(map[string]map[uint8]chWarp)
var lock = sync.RWMutex{}

// 如果needKeep为true，请不要在外部关闭XClient
func RegisterTerm(ctx context.Context, c client.XClient, ch chan *protocol.Message, sid string, platform uint8) (needKeep bool, err error) {
	lock.RLock()
	// 如果客户端channel已存在，就不要再去重复注册了
	// TODO: 这里会影响重复登录后的主动踢出已登陆客户端逻辑，主要是长轮询的场景中会反复断开重连
	// 修复方法1：msgbroker增加unregister逻辑
	if _, ok := subs[sid]; ok {
		if _, ok := subs[sid][platform]; ok {
			return false, nil
		}
	}
	lock.RUnlock()
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
		subs[sid][platform] = chWarp{
			sid:      sid,
			platform: platform,
			ch:       make(chan model.Msg),
			queue:    NewConcurrentQueue[model.Msg](),
		}
	}
	lock.Unlock()
	if regSessionReply.NeedKeep {
		// 保持连接
		select {
		case cs.ch <- &clientAndCh{
			c:  c,
			ch: ch,
		}:
			println("register success", sid, "-", platform)
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

// TODO: 因keepalive的原因，poll协议的http的ctx.Done不会马上关闭，这里可能会出现一些问题
func (cch chWarp) UnRegistry() {
	// lock.Lock()
	// defer lock.Unlock()
	// if chs, ok := subs[cch.sid]; ok {
	// 	if c, ok := chs[cch.platform]; ok {
	// 		close(c.ch)
	// 		delete(chs, cch.platform)
	// 	}
	// 	if len(chs) == 0 {
	// 		delete(subs, cch.sid)
	// 		println("unregister", cch.sid, "-", cch.platform)
	// 	}
	// }
}

func (cch chWarp) GetCh() <-chan model.Msg {
	return cch.ch
}

func RunMsgPollServer(ctx context.Context) {
	run(ctx)
}

func run(ctx context.Context) {
	for {
		if cn, ok := <-cs.ch; ok {
			go cn.ListenMsg(ctx)
		}
	}
}

// 考虑为每个client分配一个queue
func (cn *clientAndCh) ListenMsg(ctx context.Context) {
	defer cn.c.Close()
	msgData := msgbroker.BroadcastMsgWarp{}
	codec := codec.GobUtils[msgbroker.BroadcastMsgWarp]{}
	var msg *protocol.Message
	var ok bool
	for {
		msg, ok = <-cn.ch
		if !ok {
			return
		}
		if err := codec.Decode(msg.Payload, &msgData); err != nil {
			println("Decode Msg Failed: ", err)
			continue
		}
		if msgData.Action == msgbroker.Heartbeat {
			continue
		}
		lock.RLock()
		// 消息分发
		for i := range msgData.ToSessions {
			if ch, ok := subs[msgData.ToSessions[i].SessionId][msgData.ToSessions[i].Platform]; ok {
				go func() {
					timer := time.NewTimer(time.Second * 3)
					defer timer.Stop()
					select {
					case ch.ch <- msgData.MetaData:
						return
					case <-timer.C:
						// TODO: 考虑是否在此处处理心跳包，以便于踢出超时的客户端
						println("send msg timeout")
						return
					}
				}()
			}
		}
		lock.RUnlock()
	}
}
