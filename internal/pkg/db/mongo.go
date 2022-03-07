package db

import (
	"context"
	"cst/internal/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type mongodb struct {
	cfg *config.Config
	ctx context.Context
}

func NewMongoDB(cfg *config.Config, ctx context.Context) *mongodb {
	return &mongodb{
		cfg: cfg,
		ctx: ctx,
	}
}

func (m *mongodb) Dsn() string {
	config := m.cfg.MongoDb
	if config.DSN != "" {
		return config.DSN
	}

	dsn := "mongodb://"
	if len(config.Username) > 0 && len(config.Password) > 0 {
		dsn = dsn + config.Username + ":" + config.Password + "@"
	}

	dsn = dsn + config.IP + ":" + config.Port

	if len(config.Options) > 0 {
		dsn = dsn + "?" + config.Options
	}

	return dsn
}

func (m *mongodb) Open() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Dsn()))
	return client, err
}
