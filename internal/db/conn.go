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

type postgresClientOpt struct {
	port     uint16
	minConns int32
	maxConns int32
	host     string
	username string
	password string
	dbName   string
}

type pgConnOpt func(*postgresClientOpt)

func NewPostgresWithOpt(opts ...pgConnOpt) *postgresClientOpt {
	p := &postgresClientOpt{
		host:     "127.0.0.1",
		port:     5432,
		minConns: 1,
		maxConns: 10,
	}
	for i := range opts {
		opts[i](p)
	}
	return p
}

func WithHost(host string) pgConnOpt {
	return func(p *postgresClientOpt) {
		p.host = host
	}
}

func WithPort(port uint16) pgConnOpt {
	return func(p *postgresClientOpt) {
		p.port = port
	}
}

func WithUsername(username string) pgConnOpt {
	return func(p *postgresClientOpt) {
		p.username = username
	}
}

func WithPassword(password string) pgConnOpt {
	return func(p *postgresClientOpt) {
		p.password = password
	}
}

func WithDbName(dbName string) pgConnOpt {
	return func(p *postgresClientOpt) {
		p.dbName = dbName
	}
}

func WithMinConns(minConns int32) pgConnOpt {
	return func(p *postgresClientOpt) {
		p.minConns = minConns
	}
}

func WithMaxConns(maxConns int32) pgConnOpt {
	return func(p *postgresClientOpt) {
		p.maxConns = maxConns
	}
}

func (p *postgresClientOpt) GetDb() *pgxpool.Pool {
	// postgres://user:password@127.0.0.1:5432/?Timezone=Asia%2FShanghai
	conn, err := pgxpool.New(context.Background(),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/?sslmode=disable&pool_min_conns=%d&pool_max_conns=%d",
			p.username,
			p.password,
			p.host,
			p.port,
			p.minConns,
			p.maxConns,
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	if err := conn.Ping(context.TODO()); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	if _, err := conn.Exec(context.Background(), shcema); err != nil {
		log.Fatal("初始化数据表失败", err)
	}
	return conn
}
