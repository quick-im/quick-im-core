package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/services/msgid"
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
	flags = append(flags, config.GetServiceFlags(config.MsgIdSerName, 8018)...)
	flags = append(flags, config.JaegerFlags...)
	flags = append(flags, config.ConsulFlags...)
	flags = append(flags, config.LogFlags...)
	app := &cli.App{
		Name:  "msgid",
		Usage: "QuickIM消息ID模块",
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
	if err := msgid.NewServer(
		config.WithIp(conf.Services.Msgid.IP),
		config.WithPort(uint16(conf.Services.Msgid.Port)),
		config.WithOpenTracing(conf.OpenTracing),
		config.WithJeagerAgentHostPort(fmt.Sprintf("%s:%d", conf.Jaeger.Host, conf.Jaeger.Port)),
		config.WithUseConsulRegistry(true),
		config.WithConsulServers(conf.Consul.Servers...),
	).Start(ctx); err != nil {
		return err
	}
	return nil
}
