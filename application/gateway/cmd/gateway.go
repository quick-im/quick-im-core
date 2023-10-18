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
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	log.SetFlags(log.Llongfile)
	flags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "version",
			Value: "1",
		}),
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "config.yaml",
			Usage:   "配置文件路径",
		},
	}
	flags = append(flags, config.JaegerFlags...)
	flags = append(flags, config.ConsulFlags...)
	flags = append(flags, config.NatsFlags...)
	flags = append(flags, config.LogFlags...)
	app := &cli.App{
		Name:   "Gateway",
		Usage:  "QuickIM网关服务",
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
	ctx := context.Background()
	ser := server.NewApiServer(
		server.WithIp("0.0.0.0"),
		server.WithPort(8088),
		server.WithOpenTracing(true),
		server.WithJeagerServiceName("Gateway"),
		server.WithJeagerAgentHostPort("127.0.0.1:6831"),
		server.WithUseConsulRegistry(true),
		server.WithConsulServers("127.0.0.1:8500"),
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
	go msgpool.RunMsgPollServer(ctx)
	ser.InitAndStartServer(ctx)
	return nil
}
