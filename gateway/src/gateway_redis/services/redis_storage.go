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

func (rs *RedisStorage) Store(key string, value string, expiration *time.Duration) error {
	if expiration == nil {
		defaultValue := time.Duration(0)
		expiration = &defaultValue
	}

	return rs.redisClient.Set(context.Background(), key, value, *expiration).Err()
}

func (rs *RedisStorage) Read(key string) (string, error) {
	return rs.redisClient.Get(context.Background(), key).Result()
}

func (rs *RedisStorage) ReadAndDelete(key string) (string, error) {
	return rs.redisClient.GetDel(context.Background(), key).Result()
}

func (rs *RedisStorage) Delete(key string) error {
	return rs.redisClient.Del(context.Background(), key).Err()
}
