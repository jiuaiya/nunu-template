package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment string
	Name        string
	HTTP        HTTP
	MySQL       MySQL
	Redis       Redis
	Metrics     Metrics
}

type HTTP struct {
	Address           string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxBodyBytes      int64
}

type MySQL struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type Redis struct {
	Enabled  bool
	Address  string
	Password string
	DB       int
}

type Metrics struct {
	Enabled bool
}

func Load() (Config, error) {
	cfg := Config{
		Environment: envString("APP_ENV", "local"),
		Name:        envString("APP_NAME", "saas-service"),
		HTTP: HTTP{
			Address:           envString("APP_HTTP_ADDR", ":8080"),
			ReadHeaderTimeout: envDuration("APP_HTTP_READ_HEADER_TIMEOUT", 5*time.Second),
			ReadTimeout:       envDuration("APP_HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:      envDuration("APP_HTTP_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:       envDuration("APP_HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout:   envDuration("APP_HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			MaxBodyBytes:      envInt64("APP_HTTP_MAX_BODY_BYTES", 1<<20),
		},
		MySQL: MySQL{
			DSN:             envString("APP_MYSQL_DSN", "foundation:foundation@tcp(127.0.0.1:23306)/foundation?charset=utf8mb4&parseTime=true&loc=UTC"),
			MaxOpenConns:    envInt("APP_MYSQL_MAX_OPEN_CONNS", 30),
			MaxIdleConns:    envInt("APP_MYSQL_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: envDuration("APP_MYSQL_CONN_MAX_LIFETIME", 30*time.Minute),
		},
		Redis: Redis{
			Enabled:  envBool("APP_REDIS_ENABLED", true),
			Address:  envString("APP_REDIS_ADDR", "127.0.0.1:26379"),
			Password: os.Getenv("APP_REDIS_PASSWORD"),
			DB:       envInt("APP_REDIS_DB", 0),
		},
		Metrics: Metrics{Enabled: envBool("APP_METRICS_ENABLED", true)},
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("APP_NAME must not be empty")
	}
	if c.HTTP.Address == "" {
		return fmt.Errorf("APP_HTTP_ADDR must not be empty")
	}
	if c.HTTP.MaxBodyBytes <= 0 {
		return fmt.Errorf("APP_HTTP_MAX_BODY_BYTES must be positive")
	}
	if c.MySQL.DSN == "" {
		return fmt.Errorf("APP_MYSQL_DSN must not be empty")
	}
	if c.MySQL.MaxOpenConns <= 0 || c.MySQL.MaxIdleConns < 0 {
		return fmt.Errorf("invalid MySQL pool limits")
	}
	if c.MySQL.MaxIdleConns > c.MySQL.MaxOpenConns {
		return fmt.Errorf("MySQL idle connections cannot exceed open connections")
	}
	if c.Redis.Enabled && c.Redis.Address == "" {
		return fmt.Errorf("APP_REDIS_ADDR must not be empty when Redis is enabled")
	}
	return nil
}

func envString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envInt64(key string, fallback int64) int64 {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
