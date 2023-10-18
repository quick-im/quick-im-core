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
	flags = append(flags, config.RethinkDbFlags...)
	flags = append(flags, config.JaegerFlags...)
	flags = append(flags, config.ConsulFlags...)
	flags = append(flags, config.NatsFlags...)
	flags = append(flags, config.LogFlags...)
	app := &cli.App{
		Name:   "msgid",
		Usage:  "QuickIM消息ID模块",
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
	rethinkDbOpt := msgdb.NewRethinkDbWithOpt(
		msgdb.WithServer("localhost:28015"),
		msgdb.WithDb("quick-im"),
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
		config.WithIp("0.0.0.0"),
		config.WithPort(8015),
		config.WithOpenTracing(true),
		config.WithJeagerAgentHostPort("127.0.0.1:6831"),
		config.WithNatsServers("nats://127.0.0.1:4222"),
		config.WithUseConsulRegistry(true),
		config.WithConsulServers("127.0.0.1:8500"),
	).Start(ctx); err != nil {
		return err
	}
	return nil
}
