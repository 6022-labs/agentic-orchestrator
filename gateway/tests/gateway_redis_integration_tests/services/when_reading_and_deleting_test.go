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

type WhenReadingAndDeletingTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenReadingAndDeletingBeforeEach(t *testing.T) *WhenReadingAndDeletingTestingSuite {
	redisClient := container_builders.GetRedisContainerBuilder().Create()

	redisMetrics := metrics_mocks.NewMockRedisMetricsInterface(gomock.NewController(t))
	redisMetrics.EXPECT().RecordOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	sut := services.NewRedisStorage(
		redisClient,
		redisMetrics,
	)

	return &WhenReadingAndDeletingTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenReadingAndDeletingAfterEach(suite *WhenReadingAndDeletingTestingSuite) {
	container_builders.GetRedisContainerBuilder().DropConnection(suite.redisClient)
}

func TestWhenReadingAndDeleting(t *testing.T) {
	t.Parallel()

	t.Run("Given key exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the value and remove the key", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingAndDeletingBeforeEach(t)
			defer WhenReadingAndDeletingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Store(context.Background(), key, "test-value", nil)
			assert.Nil(t, err)

			value, err := suite.sut.ReadAndDelete(context.Background(), key)
			assert.Nil(t, err)
			if assert.NotNil(t, value) {
				assert.Equal(t, "test-value", *value)
			}

			value, err = suite.sut.ReadAndDelete(context.Background(), key)
			assert.Nil(t, err)
			assert.Nil(t, value)
		})
	})

	t.Run("Given key does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return no value and no error", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingAndDeletingBeforeEach(t)
			defer WhenReadingAndDeletingAfterEach(suite)

			value, err := suite.sut.ReadAndDelete(context.Background(), uuid.New().String())
			assert.Nil(t, err)
			assert.Nil(t, value)
		})
	})
}
