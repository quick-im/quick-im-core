package contant

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ContentKey string
type PgCtxType = *pgxpool.Pool
type RedisCtxType = *redis.Client
type LoggerCtxType = *zap.Logger
type NatsCtxType = *messaging.NatsWarp
type RpcxClientCtxType = *rpcx.RpcxClientWithOpt

type Sort bool

var (
	Desc Sort = true
	Asc  Sort = false
)
