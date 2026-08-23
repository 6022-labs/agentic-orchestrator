package configurations

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_redis/metrics"
	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_redis/metrics/interfaces"
	redis_services "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_redis/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_redis/settings"
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddSlackBridgeRedisConfiguration(container *dig.Container) {
	// Metrics
	if err := container.Provide(
		metrics.NewRedisMetrics,
		dig.As(new(metrics_interfaces.RedisMetricsInterface)),
	); err != nil {
		panic(err)
	}

	// Settings
	err := container.Provide(settings.NewRedisSettings)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		redis_services.NewRedisStorage,
		dig.As(new(interfaces.CacheStorageInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Provide the redis client
	err = container.Provide(func(redisSettings *settings.RedisSettings, logger *zap.Logger) *redis.Client {
		client := redis.NewClient(&redis.Options{
			Addr:     redisSettings.Host,
			Password: redisSettings.Password,
			DB:       int(redisSettings.DB),
		})

		_, err := client.Ping(context.Background()).Result()

		if err != nil {
			logger.Fatal("Error connecting to Redis: ", zap.Error(err))
		}

		return client
	})
	if err != nil {
		panic(err)
	}
}
