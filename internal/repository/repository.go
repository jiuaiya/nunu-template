package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

func New(db *gorm.DB, redisClient *redis.Client, logger *zap.Logger) *Repository {
	return &Repository{db: db, redis: redisClient, logger: logger}
}

type transactionKey struct{}

func (r *Repository) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Redis() *redis.Client {
	return r.redis
}

func (r *Repository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, transactionKey{}, tx))
	})
}
