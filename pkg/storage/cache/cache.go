package cache

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Host     string
	Port     int32
	Password string
	DB       int
}

func GetConnection(cfg *Config) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	return rdb
}
