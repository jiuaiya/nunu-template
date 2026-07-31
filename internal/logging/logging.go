package logging

import (
	"github.com/example/nunu-saas-template/internal/config"
	"go.uber.org/zap"
)

func New(cfg config.Config) (*zap.Logger, func(), error) {
	var (
		logger *zap.Logger
		err    error
	)
	if cfg.Environment == "production" || cfg.Environment == "prod" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() {
		_ = logger.Sync()
	}
	return logger, cleanup, nil
}
