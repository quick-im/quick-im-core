package contant

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/quick-im/quick-im-core/internal/cache"
	"github.com/quick-im/quick-im-core/internal/jwt"
	"github.com/quick-im/quick-im-core/internal/logger"
	"github.com/quick-im/quick-im-core/internal/messaging"
	"github.com/quick-im/quick-im-core/internal/rpcx"
	"github.com/redis/go-redis/v9"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type ContentKey string
type PgCtxType = *pgxpool.Pool
type RedisCtxType = *redis.Client
type CacheCtxType = cache.Cache
type LoggerCtxType = logger.Logger
type NatsCtxType = *messaging.NatsWarp
type RpcxClientCtxType = *rpcx.RpcxClientWithOpt
type RethinkDbCtxType = *r.Session

var JWTClaimsCtxType = &jwt.MyCustomClaims{}

type Sort bool

var (
	Desc Sort = true
	Asc  Sort = false
)
