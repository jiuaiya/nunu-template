package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/example/nunu-saas-template/internal/config"
	"github.com/example/nunu-saas-template/internal/handler"
	"github.com/example/nunu-saas-template/internal/middleware"
	"github.com/example/nunu-saas-template/internal/observability"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTP struct {
	server          *http.Server
	logger          *zap.Logger
	shutdownTimeout config.HTTP
}

func NewHTTP(cfg config.Config, logger *zap.Logger, health *handler.Health, metrics *observability.Metrics) *HTTP {
	if cfg.Environment == "production" || cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.AccessLog(logger),
		middleware.Recovery(logger),
	)
	if cfg.Metrics.Enabled {
		router.Use(metrics.Middleware())
		router.GET("/metrics", gin.WrapH(metrics.Handler()))
	}
	health.Register(router)
	router.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "pong"})
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "not_found", "message": "route not found"})
	})

	return &HTTP{
		server: &http.Server{
			Addr:              cfg.HTTP.Address,
			Handler:           http.MaxBytesHandler(router, cfg.HTTP.MaxBodyBytes),
			ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
			ReadTimeout:       cfg.HTTP.ReadTimeout,
			WriteTimeout:      cfg.HTTP.WriteTimeout,
			IdleTimeout:       cfg.HTTP.IdleTimeout,
			MaxHeaderBytes:    1 << 20,
		},
		logger:          logger,
		shutdownTimeout: cfg.HTTP,
	}
}

func (s *HTTP) Run(ctx context.Context) error {
	errorsCh := make(chan error, 1)
	go func() {
		s.logger.Info("HTTP server listening", zap.String("address", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorsCh <- fmt.Errorf("listen HTTP: %w", err)
		}
		close(errorsCh)
	}()

	select {
	case err := <-errorsCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout.ShutdownTimeout)
		defer cancel()
		if err := s.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown HTTP: %w", err)
		}
		return nil
	}
}
