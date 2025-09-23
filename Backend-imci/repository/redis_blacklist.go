// repository/redis_blacklist.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenBlacklist struct {
	client *redis.Client
}

func NewRedisTokenBlacklist(redisURL string) (*RedisTokenBlacklist, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisTokenBlacklist{client: client}, nil
}

func (r *RedisTokenBlacklist) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	err := r.client.Set(ctx, "blacklist:"+token, "1", expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	return nil
}

func (r *RedisTokenBlacklist) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	return result > 0, nil
}

func (r *RedisTokenBlacklist) Close() error {
	return r.client.Close()
}