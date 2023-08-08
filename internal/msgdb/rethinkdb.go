package msgdb

import (
	"fmt"
	"os"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type rethinkClientOpt struct {
	servers  []string
	db       string
	tables   []tableInfo
	authkey  string
	username string
	password string
}

type tableInfo struct {
	table string
	pk    string
	index []string
}

type rethinkOpt func(*rethinkClientOpt)

func NewRethinkDbWithOpt(opts ...rethinkOpt) *rethinkClientOpt {
	n := &rethinkClientOpt{
		servers: make([]string, 0),
		tables:  make([]tableInfo, 0),
	}
	for i := range opts {
		opts[i](n)
	}
	return n
}

func WithServer(server string) rethinkOpt {
	return func(nco *rethinkClientOpt) {
		nco.servers = append(nco.servers, server)
	}
}

func WithServers(servers ...string) rethinkOpt {
	return func(nco *rethinkClientOpt) {
		nco.servers = servers
	}
}

func WithDb(db string) rethinkOpt {
	return func(nco *rethinkClientOpt) {
		nco.db = db
	}
}

func WithAuthKey(authkey string) rethinkOpt {
	return func(nco *rethinkClientOpt) {
		nco.authkey = authkey
	}
}

func WithUsername(username string) rethinkOpt {
	return func(nco *rethinkClientOpt) {
		nco.username = username
	}
}

func WithPassword(password string) rethinkOpt {
	return func(nco *rethinkClientOpt) {
		nco.authkey = password
	}
}

func Table(table, pk string, indexs []string) tableInfo {
	return tableInfo{
		table: table,
		pk:    pk,
		index: indexs,
	}
}

func WithTables(table ...tableInfo) rethinkOpt {
	return func(rco *rethinkClientOpt) {
		rco.tables = append(rco.tables, table...)
	}
}

func (re *rethinkClientOpt) GetRethinkDb() *r.Session {

	session, err := r.Connect(r.ConnectOpts{
		Addresses:  re.servers, // endpoint without http
		Database:   re.db,
		InitialCap: 10,
		MaxOpen:    10,
		Username:   re.username,
		Password:   re.password,
		AuthKey:    re.authkey,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to rethinkdb: %v\n", err)
		os.Exit(1)
	}
	return session
}

func (re *rethinkClientOpt) InitDb() {
	session, err := r.Connect(r.ConnectOpts{
		Addresses: re.servers, // endpoint without http
		Database:  re.db,
		Username:  re.username,
		Password:  re.password,
		AuthKey:   re.authkey,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to rethinkdb: %v\n", err)
		os.Exit(1)
	}
	defer session.Close()
	err = r.DBList().Contains(re.db).
		Do(func(isExists r.Term) r.Term {
			return r.Branch(
				isExists,
				0,
				r.DBCreate(re.db),
			)
		}).Exec(session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rethinkdb failed to initialize the database: %s err: %v\n", re.db, err)
		os.Exit(1)
	}
	for i := range re.tables {
		err = r.DB(re.db).TableList().Contains(re.tables[i].table).
			Do(func(isExists r.Term) r.Term {
				return r.Branch(
					isExists,
					0,
					r.DB(re.db).TableCreate(re.tables[i].table, r.TableCreateOpts{
						PrimaryKey: re.tables[i].pk,
					}),
				)
			}).Exec(session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rethinkdb failed to initialize the data table: %s err: %v\n", re.tables[i].table, err)
			os.Exit(1)
		}
		for j := range re.tables[i].index {
			err = r.DB(re.db).Table(re.tables[i].table).IndexList().Contains(re.tables[i].index[j]).
				Do(func(isExists r.Term) r.Term {
					return r.Branch(
						isExists,
						0,
						r.DB(re.db).Table(re.tables[i].table).IndexCreate(re.tables[i].index[j]),
					)
				}).Exec(session)
			if err != nil {
				fmt.Fprintf(os.Stderr, "rethinkdb failed to initialize the data table: %s Index: %s err: %v\n", re.tables[i].table, re.tables[i].index[j], err)
				os.Exit(1)
			}
		}
	}
}
