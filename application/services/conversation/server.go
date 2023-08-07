package main

import (
	"context"

	"github.com/quick-im/quick-im-core/internal/contant"
	"github.com/quick-im/quick-im-core/internal/db"
	"github.com/quick-im/quick-im-core/internal/helper"
	"github.com/quick-im/quick-im-core/services/conversation"
)

func main() {
	dbOpt := db.NewPostgresWithOpt(
		db.WithHost("localhost"),
		db.WithPort(5432),
		db.WithUsername("postgres"),
		db.WithPassword("123456"),
		db.WithDbName("quickim"),
	)
	ctx := helper.InitCtx(context.Background(),
		helper.CtxOptWarp[contant.PgCtxType](contant.CTX_POSTGRES_KEY, dbOpt.GetDb()),
	)
	if err := conversation.NewServer(
		conversation.WithIp("0.0.0.0"),
		conversation.WithPort(8016),
	).Start(ctx); err != nil {
		panic(err)
	}
}
