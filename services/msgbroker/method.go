package msgbroker

import (
	"context"

	"github.com/quick-im/quick-im-core/internal/quickparam/msgbroker"
)

type BroadcastRecvFn func(context.Context, msgbroker.BroadcastArgs, *msgbroker.BroadcastReply) error

func (r *rpcxServer) BroadcastRecv(ctx context.Context) BroadcastRecvFn {
	return func(ctx context.Context, ba msgbroker.BroadcastArgs, br *msgbroker.BroadcastReply) error {
		return nil
	}
}
