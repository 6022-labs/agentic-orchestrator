package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	redisClient *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) *RedisStorage {
	return &RedisStorage{redisClient: redisClient}
}

func (rs *RedisStorage) Store(ctx context.Context, key string, value string, expiration *time.Duration) error {
	if expiration == nil {
		defaultValue := time.Duration(0)
		expiration = &defaultValue
	}

	return rs.redisClient.Set(ctx, key, value, *expiration).Err()
}

func (rs *RedisStorage) Read(ctx context.Context, key string) (string, error) {
	return rs.redisClient.Get(ctx, key).Result()
}

func (rs *RedisStorage) Delete(ctx context.Context, key string) error {
	return rs.redisClient.Del(ctx, key).Err()
}
