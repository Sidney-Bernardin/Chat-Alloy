package redis

import (
	"context"

	"github.com/Sidney-Bernardin/Chat-Alloy/internal"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	cfg    *internal.Config
	client *redis.Client
}

func New(ctx context.Context, cfg *internal.Config) (*Repository, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.REDIS_ADDR,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "cannot ping redis")
	}

	return &Repository{cfg, client}, nil
}
