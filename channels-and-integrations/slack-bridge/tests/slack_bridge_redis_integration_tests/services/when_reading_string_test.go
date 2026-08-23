package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_redis/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	redis_metrics_mocks "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/slack_bridge_redis_mocks/metrics_mocks/interfaces_mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenReadingStringTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenReadingStringBeforeEach(t *testing.T) *WhenReadingStringTestingSuite {
	redisClient := container_builders.GetRedisConnection()

	redisMetrics := redis_metrics_mocks.NewMockRedisMetricsInterface(gomock.NewController(t))
	redisMetrics.EXPECT().RecordOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	sut := services.NewRedisStorage(
		redisClient,
		redisMetrics,
	)

	return &WhenReadingStringTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenReadingStringAfterEach(suite *WhenReadingStringTestingSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenReadingString(t *testing.T) {
	t.Parallel()

	t.Run(
		"Given a key storing 'value' as string, When reading it as String, Then it returns 'value'",
		func(t *testing.T) {
			key := fmt.Sprintf("test_key_value_%s_%d", t.Name(), time.Now().UnixNano())

			suite := WhenReadingStringBeforeEach(t)
			defer WhenReadingStringAfterEach(suite)

			err := suite.redisClient.Set(context.Background(), key, "value", 0).Err()
			if err != nil {
				t.Fatalf("Failed to set key in Redis: %v", err)
			}

			result, err := suite.sut.ReadString(context.Background(), key)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "value", *result)
		},
	)

	t.Run("Given a non-existing key, When reading it as String, Then it returns nothing", func(t *testing.T) {
		key := fmt.Sprintf("non_existing_%s_%d", t.Name(), time.Now().UnixNano())

		suite := WhenReadingStringBeforeEach(t)
		defer WhenReadingStringAfterEach(suite)

		result, err := suite.sut.ReadString(context.Background(), key)

		assert.Nil(t, err)
		assert.Nil(t, result)
	})
}
