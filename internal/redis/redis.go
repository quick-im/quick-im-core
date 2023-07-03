package redis

import "github.com/redis/go-redis/v9"

func GetRedis() *redis.Client {
	rClient := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
		},
	)
	return rClient
}
