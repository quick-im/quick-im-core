package db

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

//go:embed ddl/schema.sql
var shcema string

func GetDb() *pgxpool.Pool {
	// postgres://user:password@127.0.0.1:5432/?Timezone=Asia%2FShanghai
	conn, err := pgxpool.New(context.Background(), "postgres://postgres:123456@127.0.0.1:5432/?sslmode=disable&Timezone=Asia%2FShanghai&pool_min_conns=1&pool_max_conns=10")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	if _, err := conn.Exec(context.Background(), shcema); err != nil {
		log.Fatal("初始化数据表失败", err)
	}
	return conn
}
