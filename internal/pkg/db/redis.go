package db

import (
	"context"
	"cst/internal/pkg/config"
	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	cfg *config.Config
	ctx context.Context
}

func NewRedis(cfg *config.Config, ctx context.Context) *redisCache {
	return &redisCache{
		cfg: cfg,
		ctx: ctx,
	}
}

// Opening a database and save the reference to `Database` struct.
func (m *redisCache) Connect() *redis.Client {
	config := m.cfg.Redis
	return redis.NewClient(&redis.Options{
		Addr:     config.Hostname + ":" + config.Port,
		Password: config.Password,
		DB:       config.Database,
	})
}
