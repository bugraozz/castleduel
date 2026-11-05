package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadFromUsesDefaults(t *testing.T) {
	values := map[string]string{
		"DATABASE_URL": "postgres://app:secret@localhost:5432/kale?sslmode=disable",
		"REDIS_URL":    "redis://localhost:6379/0",
	}

	got, err := LoadFrom(func(name string) string { return values[name] })
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}
	if got.Environment != "development" || got.APIAddress != ":8080" {
		t.Fatalf("unexpected defaults: environment=%q address=%q", got.Environment, got.APIAddress)
	}
	if got.ShutdownTimeout != 10*time.Second {
		t.Fatalf("ShutdownTimeout = %s, want 10s", got.ShutdownTimeout)
	}
}

func TestLoadFromRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]string
		want   string
	}{
		{
			name:   "missing database URL",
			values: map[string]string{"REDIS_URL": "redis://localhost:6379/0"},
			want:   "DATABASE_URL is required",
		},
		{
			name: "invalid API address",
			values: map[string]string{
				"API_ADDR":     "not-an-address",
				"DATABASE_URL": "postgres://app:secret@localhost:5432/kale",
				"REDIS_URL":    "redis://localhost:6379/0",
			},
			want: "API_ADDR must be a valid host:port address",
		},
		{
			name: "invalid Redis URL",
			values: map[string]string{
				"DATABASE_URL": "postgres://app:secret@localhost:5432/kale",
				"REDIS_URL":    "not-a-url",
			},
			want: "REDIS_URL is invalid",
		},
		{
			name: "invalid shutdown timeout",
			values: map[string]string{
				"DATABASE_URL":     "postgres://app:secret@localhost:5432/kale",
				"REDIS_URL":        "redis://localhost:6379/0",
				"SHUTDOWN_TIMEOUT": "0s",
			},
			want: "SHUTDOWN_TIMEOUT must be a positive duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadFrom(func(name string) string { return tt.values[name] })
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("LoadFrom() error = %v, want containing %q", err, tt.want)
			}
		})
	}
}
