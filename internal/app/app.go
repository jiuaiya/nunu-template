package app

import (
	"context"
	"github.com/example/nunu-saas-template/internal/server"
	"go.uber.org/zap"
)

type App struct {
	server *server.HTTP
	logger *zap.Logger
}

func New(httpServer *server.HTTP, logger *zap.Logger) *App {
	return &App{server: httpServer, logger: logger}
}

func (a *App) Run(ctx context.Context) error {
	return a.server.Run(ctx)
}
