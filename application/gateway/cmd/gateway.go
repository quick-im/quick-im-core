package main

import (
	"context"
	"log"
	"os"

	"github.com/quick-im/quick-im-core/application/gateway/gateway/msgpool"
	"github.com/quick-im/quick-im-core/application/gateway/gateway/server"
	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/services/conversation"
	"github.com/quick-im/quick-im-core/services/msgbroker"
	"github.com/quick-im/quick-im-core/services/msghub"
	"github.com/quick-im/quick-im-core/services/msgid"
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
	flags = append(flags, config.GetServiceFlags(config.GatewayName, 8088)...)
	flags = append(flags, config.GatewayIpwriteFlags...)
	flags = append(flags, config.JaegerFlags...)
	flags = append(flags, config.ConsulFlags...)
	flags = append(flags, config.NatsFlags...)
	flags = append(flags, config.LogFlags...)
	app := &cli.App{
		Name:  "Gateway",
		Usage: "QuickIM网关服务",
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
	ser := server.NewApiServer(
		server.WithIp(conf.Gateway.IP),
		server.WithPort(uint16(conf.Gateway.Port)),
		server.WithOpenTracing(conf.OpenTracing),
		server.WithJeagerServiceName("Gateway"),
		server.WithJeagerAgentHostPort(conf.Jaeger.Host),
		server.WithUseConsulRegistry(true),
		server.WithConsulServers(conf.Consul.Servers...),
	)
	persistence := ser.InitDepServices(persistence.SERVER_NAME)
	msgbroker := ser.InitDepServices(msgbroker.SERVER_NAME)
	msghub := ser.InitDepServices(msghub.SERVER_NAME)
	msgid := ser.InitDepServices(msgid.SERVER_NAME)
	conversation := ser.InitDepServices(conversation.SERVER_NAME)
	defer func() {
		_ = persistence.Close()
		_ = msgbroker.Close()
		_ = msghub.Close()
		_ = msgid.Close()
		_ = conversation.Close()
	}()
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_PERSISTENCE, persistence)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_MSGBORKER, msgbroker)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_MSGHUB, msghub)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_MSGID, msgid)
	ctx = context.WithValue(ctx, contant.CTX_SERVICE_CONVERSATION, conversation)
	ctx = context.WithValue(ctx, contant.CTX_IP_WHITELIST_KEY, conf.Gateway.IPWrite)
	go msgpool.RunMsgPollServer(ctx)
	ser.InitAndStartServer(ctx)
	return nil
}
