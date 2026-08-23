package services

import (
	"context"
	"errors"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_redis/metrics/interfaces"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	redisClient  *redis.Client
	redisMetrics metrics_interfaces.RedisMetricsInterface
}

func NewRedisStorage(
	redisClient *redis.Client,
	redisMetrics metrics_interfaces.RedisMetricsInterface,
) *RedisStorage {
	return &RedisStorage{redisClient: redisClient, redisMetrics: redisMetrics}
}

func (rs *RedisStorage) Store(ctx context.Context, key string, value string, expiration *time.Duration) error {
	if expiration == nil {
		defaultValue := time.Duration(0)
		expiration = &defaultValue
	}

	start := time.Now()
	err := rs.redisClient.Set(ctx, key, value, *expiration).Err()
	rs.redisMetrics.RecordOperation(ctx, "SET", time.Since(start), err)
	return err
}

func (rs *RedisStorage) ReadAndDelete(ctx context.Context, key string) (*string, error) {
	start := time.Now()
	val, err := rs.redisClient.GetDel(ctx, key).Result()
	// redis.Nil is a miss, not an operational failure.
	rs.redisMetrics.RecordOperation(ctx, "GETDEL", time.Since(start), rs.operationError(err))
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return &val, nil
}

// operationError treats a missing key as a miss rather than an operational failure.
func (rs *RedisStorage) operationError(err error) error {
	if err == nil || errors.Is(err, redis.Nil) {
		return nil
	}
	return err
}
