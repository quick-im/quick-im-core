package contant

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ContentKey string
type PgCtxType = *pgxpool.Pool
type RedisCtxType = *redis.Client
type LoggerCtxType = *zap.Logger
