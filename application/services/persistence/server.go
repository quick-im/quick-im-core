package main

import (
	"context"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/internal/msgdb"
	"github.com/quick-im/quick-im-core/services/persistence"
)

func main() {
	dbOpt := db.NewPostgresWithOpt(
		db.WithHost("localhost"),
		db.WithPort(5432),
		db.WithUsername("postgres"),
		db.WithPassword("123456"),
		db.WithDbName("quickim"),
	)
	rethinkDbOpt := msgdb.NewRethinkDbWithOpt(
		msgdb.WithServer("localhost:28015"),
		msgdb.WithDb("quick-im"),
		msgdb.WithTables(
			msgdb.Table("msg", "msg_id", []string{"msg_id1"}),
		),
	)
	rethinkDbOpt.InitDb()
	ctx := helper.InitCtx(context.Background(),
		helper.CtxOptWarp[contant.PgCtxType](contant.CTX_POSTGRES_KEY, dbOpt.GetDb()),
		helper.CtxOptWarp[contant.RethinkDbCtxType](contant.CTX_POSTGRES_KEY, rethinkDbOpt.GetRethinkDb()),
	)
	if err := persistence.NewServer(
		persistence.WithIp("0.0.0.0"),
		persistence.WithPort(8015),
	).Start(ctx); err != nil {
		panic(err)
	}
}
