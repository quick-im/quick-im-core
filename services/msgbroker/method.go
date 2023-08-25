package msgbroker

import (
	"context"
	"net"

	"github.com/quick-im/quick-im-core/internal/quickparam/msgbroker"
	"github.com/smallnest/rpcx/server"
)

type BroadcastRecvFn func(context.Context, msgbroker.BroadcastArgs, *msgbroker.BroadcastReply) error

func (r *rpcxServer) BroadcastRecv(ctx context.Context) BroadcastRecvFn {
	return func(ctx context.Context, ba msgbroker.BroadcastArgs, br *msgbroker.BroadcastReply) error {
		return nil
	}
}

type RegisterSessionFn func(context.Context, msgbroker.RegisterSessionInfo, msgbroker.RegisterSessionReply) error

// 将接入层的Session注册到broker，用于消息分发，这里可以优化一下，在注册之前踢掉其他同Session同Platform的客户端
func (r *rpcxServer) RegisterSession(ctx context.Context) RegisterSessionFn {
	return func(ctx context.Context, args msgbroker.RegisterSessionInfo, reply msgbroker.RegisterSessionReply) error {
		clientConn := ctx.Value(server.RemoteConnContextKey).(net.Conn)
		args.Conn = clientConn
		r.connList[args.SessionId] = args
		return nil
	}
}
