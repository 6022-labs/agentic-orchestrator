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

func (rs *RedisStorage) Store(ctx context.Context, key string, value any, expiration *time.Duration) error {
	if expiration == nil {
		defaultValue := time.Duration(0)
		expiration = &defaultValue
	}

	return rs.redisClient.Set(ctx, key, value, *expiration).Err()
}

func (rs *RedisStorage) ReadBool(ctx context.Context, key string) (*bool, error) {
	val, err := rs.redisClient.Get(ctx, key).Bool()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return &val, nil
}

func (rs *RedisStorage) ReadString(ctx context.Context, key string) (*string, error) {
	val, err := rs.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return &val, nil
}
