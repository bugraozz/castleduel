package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"castleduello.local/backend/internal/config"
	"castleduello.local/backend/internal/runtime"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("worker stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	database, err := runtime.OpenDatabase(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer database.Close()

	cache, err := runtime.OpenRedis(cfg.RedisURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := cache.Close(); err != nil {
			logger.Error("failed to close Redis client", "error", err)
		}
	}()

	if err := runtime.CheckDependencies(context.Background(), database, cache); err != nil {
		return err
	}
	logger.Info("worker ready")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	logger.Info("worker stopped gracefully")
	return nil
}
