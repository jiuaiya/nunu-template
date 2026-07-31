//go:build wireinject

package wire

import (
	"context"

	"github.com/example/nunu-saas-template/internal/app"
	"github.com/example/nunu-saas-template/internal/config"
	"github.com/example/nunu-saas-template/internal/handler"
	"github.com/example/nunu-saas-template/internal/logging"
	"github.com/example/nunu-saas-template/internal/observability"
	"github.com/example/nunu-saas-template/internal/platform/cache"
	"github.com/example/nunu-saas-template/internal/platform/database"
	"github.com/example/nunu-saas-template/internal/server"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func newHealth(mysql *database.MySQL, redisClient *redis.Client) *handler.Health {
	checks := []handler.DependencyCheck{
		{Name: "mysql", Check: mysql.SQL.PingContext},
	}
	if redisClient != nil {
		checks = append(checks, handler.DependencyCheck{
			Name: "redis",
			Check: func(ctx context.Context) error {
				return redisClient.Ping(ctx).Err()
			},
		})
	}
	return handler.NewHealth(checks...)
}

func newMetrics(cfg config.Config) *observability.Metrics {
	return observability.NewMetrics(cfg.Name)
}

func New(context.Context, config.Config) (*app.App, func(), error) {
	wire.Build(
		logging.New,
		database.OpenMySQL,
		cache.OpenRedis,
		newHealth,
		newMetrics,
		server.NewHTTP,
		app.New,
	)
	return nil, nil, nil
}
