package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/example/nunu-saas-template/cmd/server/wire"
	"github.com/example/nunu-saas-template/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, cleanup, err := wire.New(ctx, cfg)
	if err != nil {
		log.Fatalf("initialize application: %v", err)
	}
	defer cleanup()

	if err := application.Run(ctx); err != nil {
		log.Fatalf("run application: %v", err)
	}
}
