package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(addr string) *Redis {
	if addr == "" {
		return &Redis{}
	}
	return &Redis{Client: redis.NewClient(&redis.Options{Addr: addr})}
}

func (r *Redis) Ping(ctx context.Context) error {
	if r.Client == nil {
		return nil
	}
	return r.Client.Ping(ctx).Err()
}

func (r *Redis) Close() error {
	if r.Client == nil {
		return nil
	}
	return r.Client.Close()
}
