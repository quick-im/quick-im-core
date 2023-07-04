package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisClientOpt struct {
	host     string
	port     uint16
	username string
	password string
}

type redisOpt func(*redisClientOpt)

func NewRedisWithOpt(opts ...redisOpt) *redisClientOpt {
	r := &redisClientOpt{
		host: "127.0.0.1",
		port: 6379,
	}
	for i := range opts {
		opts[i](r)
	}
	return r
}

func WithHost(host string) redisOpt {
	return func(rco *redisClientOpt) {
		rco.host = host
	}
}

func WithPost(port uint16) redisOpt {
	return func(rco *redisClientOpt) {
		rco.port = port
	}
}

func WithUsername(username string) redisOpt {
	return func(rco *redisClientOpt) {
		rco.username = username
	}
}

func WithPassword(password string) redisOpt {
	return func(rco *redisClientOpt) {
		rco.password = password
	}
}

func (r *redisClientOpt) GetRedis() *redis.Client {
	rClient := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", r.host, r.port),
			Username: r.username,
			Password: r.password,
		},
	)
	return rClient
}
