package main

import (
	"context"
	"log"
	"os"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/msgdb"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"github.com/quick-im/quick-im-core/services/persistence"
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
	flags = append(flags, config.GetServiceFlags(config.PersistenceSerName, 8015)...)
	flags = append(flags, config.RethinkDbFlags...)
	flags = append(flags, config.JaegerFlags...)
	flags = append(flags, config.ConsulFlags...)
	flags = append(flags, config.NatsFlags...)
	flags = append(flags, config.LogFlags...)
	app := &cli.App{
		Name:  "persistence",
		Usage: "QuickIM消息持久化模块",
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
	rethinkDbOpt := msgdb.NewRethinkDbWithOpt(
		msgdb.WithServers(conf.RethinkDb.Servers...),
		msgdb.WithDb(conf.RethinkDb.Db),
		msgdb.WithTables(
			// msgdb.Table("msg", "msg_id", []string{"msg_id1"}),
			// 二者选其一即可，可手动指定索引，也可通过model来自动创建
			msgdb.Model(model.Msg{}),
		),
	)
	rethinkDbOpt.InitDb()
	ctx := helper.InitCtx(context.Background(),
		helper.CtxOptWarp[contant.RethinkDbCtxType](contant.CTX_RETHINK_DB_KEY, rethinkDbOpt.GetRethinkDb()),
	)
	if err := persistence.NewServer(
		config.WithIp(conf.Services.Persistence.IP),
		config.WithPort(uint16(conf.Services.Persistence.Port)),
		config.WithOpenTracing(conf.OpenTracing),
		config.WithJeagerAgentHostPort(conf.Jaeger.Host),
		config.WithNatsServers(conf.Nats.Servers...),
		config.WithUseConsulRegistry(true),
		config.WithConsulServers(conf.Consul.Servers...),
	).Start(ctx); err != nil {
		return err
	}
	return nil
}
