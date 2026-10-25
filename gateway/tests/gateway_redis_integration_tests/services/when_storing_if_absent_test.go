package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_redis/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	metrics_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_redis_mocks/metrics_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenStoringIfAbsentTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenStoringIfAbsentBeforeEach(t *testing.T) *WhenStoringIfAbsentTestingSuite {
	redisClient := container_builders.GetRedisConnection()

	redisMetrics := metrics_mocks.NewMockRedisMetricsInterface(gomock.NewController(t))
	redisMetrics.EXPECT().RecordOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	return &WhenStoringIfAbsentTestingSuite{
		redisClient: redisClient,
		sut:         services.NewRedisStorage(redisClient, redisMetrics),
	}
}

func WhenStoringIfAbsentAfterEach(suite *WhenStoringIfAbsentTestingSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenStoringIfAbsent(t *testing.T) {
	t.Parallel()

	t.Run("Given key does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store it and report it", func(t *testing.T) {
			t.Parallel()
			suite := WhenStoringIfAbsentBeforeEach(t)
			defer WhenStoringIfAbsentAfterEach(suite)

			stored, err := suite.sut.StoreIfAbsent(context.Background(), uuid.New().String(), "1", time.Minute)

			assert.Nil(t, err)
			assert.True(t, stored)
		})
	})

	t.Run("Given key already exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not store it again", func(t *testing.T) {
			t.Parallel()
			suite := WhenStoringIfAbsentBeforeEach(t)
			defer WhenStoringIfAbsentAfterEach(suite)
			key := uuid.New().String()
			_, err := suite.sut.StoreIfAbsent(context.Background(), key, "1", time.Minute)
			assert.Nil(t, err)

			stored, err := suite.sut.StoreIfAbsent(context.Background(), key, "2", time.Minute)

			assert.Nil(t, err)
			assert.False(t, stored)
		})
	})
}
