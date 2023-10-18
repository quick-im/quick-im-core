package main

import (
	"context"
	"log"
	"os"

	"github.com/quick-im/quick-im-core/internal/cache/redis"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/services/conversation"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	flags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name: "version",
		}),
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "config.yaml",
			Usage:   "配置文件路径",
		},
	}
	flags = append(flags, config.BaseFlags...)
	app := &cli.App{
		Name:   "conversation",
		Usage:  "QuickIM会话管理模块",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
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
		config.WithIp("0.0.0.0"),
		config.WithPort(8016),
		config.WithOpenTracing(true),
		config.WithJeagerAgentHostPort("127.0.0.1:6831"),
		config.WithUseConsulRegistry(true),
		config.WithConsulServers("127.0.0.1:8500"),
	).Start(ctx); err != nil {
		return err
	}
	return nil
}
