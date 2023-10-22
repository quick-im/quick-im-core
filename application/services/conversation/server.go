package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/quick-im/quick-im-core/internal/cache/redis"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/services/conversation"
	"github.com/urfave/cli/v2"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "config.toml",
			Usage:   "配置文件路径",
		},
	}
	flags = append(flags, config.GetServiceFlags(config.ConversationSerName, 8016)...)
	flags = append(flags, config.JaegerFlags...)
	flags = append(flags, config.PgFlags...)
	flags = append(flags, config.ConsulFlags...)
	flags = append(flags, config.RedisFlags...)
	flags = append(flags, config.NatsFlags...)
	flags = append(flags, config.LogFlags...)
	app := &cli.App{
		Name:  "conversation",
		Usage: "QuickIM会话管理模块",
		Flags: flags,
		Action: func(ctx *cli.Context) error {
			err := run(ctx)
			return err
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args *cli.Context) error {
	conf := config.MergeConf(args)
	// fmt.Printf("%#v\n", conf)
	dbOpt := db.NewPostgresWithOpt(
		db.WithHost(conf.Postgres.Host),
		db.WithPort(uint16(conf.Postgres.Port)),
		db.WithUsername(conf.Postgres.Username),
		db.WithPassword(conf.Postgres.Password),
		db.WithDbName(conf.Postgres.Db),
	)
	cacheOpt := redis.NewRedisWithOpt(
		redis.WithHost(conf.Redis.Host),
		redis.WithPost(uint16(conf.Redis.Port)),
		redis.WithUsername(conf.Redis.Username),
		redis.WithPassword(conf.Redis.Password),
	)
	ctx := helper.InitCtx(context.Background(),
		helper.CtxOptWarp[contant.PgCtxType](contant.CTX_POSTGRES_KEY, dbOpt.GetDb()),
		helper.CtxOptWarp[contant.CacheCtxType](contant.CTX_CACHE_DB_KEY, cacheOpt.GetRedis()),
	)
	if err := conversation.NewServer(
		config.WithIp(conf.Services.Conversation.IP),
		config.WithPort(uint16(conf.Services.Conversation.Port)),
		config.WithOpenTracing(conf.OpenTracing),
		config.WithJeagerAgentHostPort(fmt.Sprintf("%s:%d", conf.Jaeger.Host, conf.Jaeger.Port)),
		config.WithUseConsulRegistry(true),
		config.WithConsulServers(conf.Consul.Servers...),
	).Start(ctx); err != nil {
		return err
	}
	return nil
}
