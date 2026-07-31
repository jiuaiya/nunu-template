package cache

import (
	"context"
	"fmt"

	"github.com/example/nunu-saas-template/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func OpenRedis(ctx context.Context, appConfig config.Config, logger *zap.Logger) (*redis.Client, func(), error) {
	cfg := appConfig.Redis
	if !cfg.Enabled {
		logger.Info("Redis disabled")
		return nil, func() {}, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, func() {}, fmt.Errorf("ping Redis: %w", err)
	}
	logger.Info("Redis connection established")
	cleanup := func() {
		_ = client.Close()
	}
	return client, cleanup, nil
}
