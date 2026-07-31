package config

import (
	"testing"
	"time"
)

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("APP_NAME", "billing")
	t.Setenv("APP_HTTP_ADDR", "127.0.0.1:9090")
	t.Setenv("APP_HTTP_MAX_BODY_BYTES", "2048")
	t.Setenv("APP_MYSQL_DSN", "user:pass@tcp(localhost:3306)/test")
	t.Setenv("APP_REDIS_ENABLED", "false")
	t.Setenv("APP_HTTP_SHUTDOWN_TIMEOUT", "3s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Name != "billing" || cfg.HTTP.Address != "127.0.0.1:9090" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.HTTP.MaxBodyBytes != 2048 || cfg.HTTP.ShutdownTimeout != 3*time.Second {
		t.Fatalf("unexpected HTTP config: %+v", cfg.HTTP)
	}
	if cfg.Redis.Enabled {
		t.Fatal("Redis should be disabled")
	}
}

func TestValidateRejectsInvalidPool(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	cfg.MySQL.MaxOpenConns = 1
	cfg.MySQL.MaxIdleConns = 2
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() expected an error")
	}
}
