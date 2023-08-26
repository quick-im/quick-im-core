package main

import (
	"context"

	"github.com/quick-im/quick-im-core/internal/cache/redis"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/services/conversation"
)

func main() {
	dbOpt := db.NewPostgresWithOpt(
		db.WithHost("localhost"),
		db.WithPort(5432),
		db.WithUsername("postgres"),
		db.WithPassword("123456"),
		db.WithDbName("quickim"),
	)
	cacheOpt := redis.NewRedisWithOpt(
		redis.WithHost("127.0.0.1"),
		redis.WithPost(6379),
	)
	ctx := helper.InitCtx(context.Background(),
		helper.CtxOptWarp[contant.PgCtxType](contant.CTX_POSTGRES_KEY, dbOpt.GetDb()),
		helper.CtxOptWarp[contant.CacheCtxType](contant.CTX_CACHE_DB_KEY, cacheOpt.GetRedis()),
	)
	if err := conversation.NewServer(
		conversation.WithIp("0.0.0.0"),
		conversation.WithPort(8016),
		conversation.WithOpenTracing(true),
		conversation.WithJeagerServiceName(conversation.SERVER_NAME),
		conversation.WithJeagerAgentHostPort("127.0.0.1:6831"),
		conversation.WithUseConsulRegistry(true),
		conversation.WithConsulServers("127.0.0.1:8500"),
	).Start(ctx); err != nil {
		panic(err)
	}
}
