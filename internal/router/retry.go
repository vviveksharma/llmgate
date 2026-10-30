package router

import (
	"context"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	Backoff     time.Duration
}

func WithRetry[T any](ctx context.Context, cfg RetryConfig, fn func(context.Context) (T, error)) (T, error) {
	var zero T
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}
	if cfg.Backoff <= 0 {
		cfg.Backoff = 200 * time.Millisecond
	}
	var err error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		var v T
		v, err = fn(ctx)
		if err == nil {
			return v, nil
		}
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-time.After(cfg.Backoff):
		}
	}
	return zero, err
}
