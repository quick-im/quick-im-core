package persistence

import (
	"fmt"
	"os"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type rethinkClientOpt struct {
	servers  []string
	db       string
	authkey  string
	username string
	password string
}

type rethinkOpt func(*rethinkClientOpt)

func NewNatsWithOpt(opts ...rethinkOpt) *rethinkClientOpt {
	n := &rethinkClientOpt{
		servers: make([]string, 0),
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
	err = r.DB("test").TableList().Contains("oldcat").
		Do(func(isExists r.Term) r.Term {
			return r.Branch(
				isExists,
				0,
				r.DB("test").TableCreate("oldcat", r.TableCreateOpts{
					PrimaryKey: "msg_id",
				}),
			)
		}).Exec(session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rethinkdb failed to initialize the data table: %v\n", err)
		os.Exit(1)
	}
	return session
}
