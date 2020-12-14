package config

import (
	"restaurant-assistant/pkg/storage/cache"
	"restaurant-assistant/pkg/storage/database"
)

const (
	ServiceName = "restaurant-assistant"

	// TODO: move to DB config
	DefaultDbHost     = "localhost"
	DefaultDbPort     = 5432
	DefaultDbUser     = "postgres"
	DefaultDbPassword = "postgres"
	DefaultDbName     = "restaurant_assistant"

	// TODO: move it
	DefaultCacheHost     = "localhost"
	DefaultCachePort     = 6379
	DefaultCachePassword = ""
	DefaultCacheDB       = 0
)

type Config struct {
	DB    *database.Config
	Cache *cache.Config
}

func Init() *Config {
	cfg := &Config{
		DB: &database.Config{
			Host:     DefaultDbHost,
			User:     DefaultDbUser,
			Password: DefaultDbPassword,
			DbName:   DefaultDbName,
			Port:     DefaultDbPort,
		},
		Cache: &cache.Config{
			Host:     DefaultCacheHost,
			Port:     DefaultCachePort,
			Password: DefaultCachePassword,
			DB:       DefaultCacheDB,
		},
	}

	return cfg
}
