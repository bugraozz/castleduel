package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"castleduello.local/backend/internal/config"
	"castleduello.local/backend/internal/health"
	"castleduello.local/backend/internal/runtime"
)

func main() {
	slog.SetDefault(newLogger("production"))
	if err := run(); err != nil {
		slog.Error("API stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logger := newLogger(cfg.Environment)
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

	mux := http.NewServeMux()
	mux.Handle("/", health.NewHandler(logger,
		health.Check{Name: "database", Run: func(ctx context.Context) error { return database.Ping(ctx) }},
		health.Check{Name: "redis", Run: func(ctx context.Context) error { return cache.Ping(ctx).Err() }},
	))
	server := &http.Server{
		Addr:              cfg.APIAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("API listening", "address", cfg.APIAddress)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		logger.Info("API stopped gracefully")
		return nil
	}
}

func newLogger(environment string) *slog.Logger {
	level := slog.LevelInfo
	if environment == "development" {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
