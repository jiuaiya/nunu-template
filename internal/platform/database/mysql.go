package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/example/nunu-saas-template/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

type MySQL struct {
	DB  *gorm.DB
	SQL *sql.DB
}

func OpenMySQL(ctx context.Context, appConfig config.Config, logger *zap.Logger) (*MySQL, func(), error) {
	cfg := appConfig.MySQL
	logLevel := gormlog.Warn
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: gormlog.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, func() {}, fmt.Errorf("open MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, func() {}, fmt.Errorf("get MySQL pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, func() {}, fmt.Errorf("ping MySQL: %w", err)
	}
	logger.Info("MySQL connection established")
	cleanup := func() {
		_ = sqlDB.Close()
	}
	return &MySQL{DB: db, SQL: sqlDB}, cleanup, nil
}
