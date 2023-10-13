package main

import (
	"context"

	"github.com/quick-im/quick-im-core/internal/config"
	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/msgdb"
	"github.com/quick-im/quick-im-core/internal/msgdb/model"
	"github.com/quick-im/quick-im-core/services/persistence"
)

func main() {
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
		panic(err)
	}
}
