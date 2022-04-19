package datastore

import (
	"github.com/go-redis/redis/v8"
	"github.com/jelliflix/jelliflix/infrastructure/config"
)

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.Addr,
		Password: config.Cfg.Redis.Pass,
	})
}
