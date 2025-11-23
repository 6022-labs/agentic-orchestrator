package configurations

import (
	"context"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services/interfaces"
	redis_services "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_redis/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_redis/settings"
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func AddOrchestratorRedisConfiguration(container *dig.Container) {
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
