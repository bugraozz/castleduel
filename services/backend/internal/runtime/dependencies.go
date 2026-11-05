package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func CheckDependencies(ctx context.Context, database *pgxpool.Pool, cache *redis.Client) error {
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := database.Ping(checkCtx); err != nil {
		return fmt.Errorf("database ping: %w", err)
	}
	if err := cache.Ping(checkCtx).Err(); err != nil {
		return fmt.Errorf("Redis ping: %w", err)
	}
	return nil
}
