package services_test

import (
	"context"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_redis/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	metrics_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_redis_mocks/metrics_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenReadingTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenReadingBeforeEach(t *testing.T) *WhenReadingTestingSuite {
	redisClient := container_builders.GetRedisContainerBuilder().Create()

	redisMetrics := metrics_mocks.NewMockRedisMetricsInterface(gomock.NewController(t))
	redisMetrics.EXPECT().RecordOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	sut := services.NewRedisStorage(
		redisClient,
		redisMetrics,
	)

	return &WhenReadingTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenReadingAfterEach(suite *WhenReadingTestingSuite) {
	container_builders.GetRedisContainerBuilder().DropConnection(suite.redisClient)
}

func TestWhenReading(t *testing.T) {
	t.Parallel()

	t.Run("Given key exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the value and keep the key", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingBeforeEach(t)
			defer WhenReadingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Store(context.Background(), key, "test-value", nil)
			assert.Nil(t, err)

			value, err := suite.sut.Read(context.Background(), key)
			assert.Nil(t, err)
			if assert.NotNil(t, value) {
				assert.Equal(t, "test-value", *value)
			}

			value, err = suite.sut.Read(context.Background(), key)
			assert.Nil(t, err)
			if assert.NotNil(t, value) {
				assert.Equal(t, "test-value", *value)
			}
		})
	})

	t.Run("Given key does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return no value and no error", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingBeforeEach(t)
			defer WhenReadingAfterEach(suite)

			value, err := suite.sut.Read(context.Background(), uuid.New().String())
			assert.Nil(t, err)
			assert.Nil(t, value)
		})
	})
}
