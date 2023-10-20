package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/services/msgbroker"
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
	flags = append(flags, config.GetServiceFlags(config.MsgbrokerSerName, 8017)...)
	flags = append(flags, config.JaegerFlags...)
	flags = append(flags, config.ConsulFlags...)
	flags = append(flags, config.NatsFlags...)
	flags = append(flags, config.LogFlags...)
	app := &cli.App{
		Name:  "msgbroker",
		Usage: "QuickIM消息分发模块",
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
	ctx := context.Background()
	if err := msgbroker.NewServer(
		config.WithIp(conf.Services.Msgbroker.IP),
		config.WithPort(uint16(conf.Services.Msghub.Port)),
		config.WithOpenTracing(conf.OpenTracing),
		config.WithJeagerAgentHostPort(fmt.Sprintf("%s:%d", conf.Jaeger.Host, conf.Jaeger.Port)),
		config.WithNatsServers(conf.Nats.Servers...),
		config.WithUseConsulRegistry(true),
		config.WithConsulServers(conf.Consul.Servers...),
	).Start(ctx); err != nil {
		return err
	}
	return nil
}
