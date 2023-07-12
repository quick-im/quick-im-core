package contant

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ContentKey string
type PgCtxType = *pgxpool.Pool
type RedisCtxType = *redis.Client
type LoggerCtxType = *zap.Logger
type NatsCtxType = *nats.Conn

type Sort bool

var (
	Desc Sort = true
	Asc  Sort = false
)

type msgGroupTopic string

var (
	PersistenceGroup msgGroupTopic = "msg.to.persistence"
)
