package services

import (
	"context"
	"errors"
	"time"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_redis/metrics/interfaces"
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

func (rs *RedisStorage) Read(ctx context.Context, key string) (string, error) {
	start := time.Now()
	value, err := rs.redisClient.Get(ctx, key).Result()
	// redis.Nil is a miss, not an operational failure.
	rs.redisMetrics.RecordOperation(ctx, "GET", time.Since(start), rs.operationError(err))
	return value, err
}

func (rs *RedisStorage) Delete(ctx context.Context, key string) error {
	start := time.Now()
	err := rs.redisClient.Del(ctx, key).Err()
	rs.redisMetrics.RecordOperation(ctx, "DEL", time.Since(start), err)
	return err
}

// operationError treats a missing key as a miss rather than an operational failure.
func (rs *RedisStorage) operationError(err error) error {
	if err == nil || errors.Is(err, redis.Nil) {
		return nil
	}
	return err
}
