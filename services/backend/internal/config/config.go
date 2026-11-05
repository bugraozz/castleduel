package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Environment     string
	APIAddress      string
	DatabaseURL     string
	RedisURL        string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	return LoadFrom(lookupEnvironment)
}

func LoadFrom(getenv func(string) string) (Config, error) {
	environment := valueOr(getenv("APP_ENV"), "development")
	switch environment {
	case "development", "test", "staging", "production":
	default:
		return Config{}, errors.New("APP_ENV must be development, test, staging, or production")
	}

	apiAddress := valueOr(getenv("API_ADDR"), ":8080")
	if _, _, err := net.SplitHostPort(apiAddress); err != nil {
		return Config{}, errors.New("API_ADDR must be a valid host:port address")
	}

	databaseURL, err := required(getenv, "DATABASE_URL")
	if err != nil {
		return Config{}, err
	}
	if _, err := pgxpool.ParseConfig(databaseURL); err != nil {
		return Config{}, errors.New("DATABASE_URL is invalid")
	}

	redisURL, err := required(getenv, "REDIS_URL")
	if err != nil {
		return Config{}, err
	}
	if _, err := redis.ParseURL(redisURL); err != nil {
		return Config{}, errors.New("REDIS_URL is invalid")
	}

	shutdownTimeout := 10 * time.Second
	if raw := strings.TrimSpace(getenv("SHUTDOWN_TIMEOUT")); raw != "" {
		shutdownTimeout, err = time.ParseDuration(raw)
		if err != nil || shutdownTimeout <= 0 {
			return Config{}, errors.New("SHUTDOWN_TIMEOUT must be a positive duration")
		}
	}

	return Config{
		Environment:     environment,
		APIAddress:      apiAddress,
		DatabaseURL:     databaseURL,
		RedisURL:        redisURL,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

func required(getenv func(string) string, name string) (string, error) {
	value := strings.TrimSpace(getenv(name))
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func lookupEnvironment(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}
