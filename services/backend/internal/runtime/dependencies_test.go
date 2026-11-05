package runtime

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestCheckDependencies(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	redisURL := os.Getenv("REDIS_URL")
	if databaseURL == "" || redisURL == "" {
		t.Skip("DATABASE_URL and REDIS_URL are required for the integration test")
	}

	database, err := OpenDatabase(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("OpenDatabase() error = %v", err)
	}
	t.Cleanup(database.Close)

	cache, err := OpenRedis(redisURL)
	if err != nil {
		t.Fatalf("OpenRedis() error = %v", err)
	}
	t.Cleanup(func() {
		if err := cache.Close(); err != nil {
			t.Errorf("Redis Close() error = %v", err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := CheckDependencies(ctx, database, cache); err != nil {
		t.Fatalf("CheckDependencies() error = %v", err)
	}
}
