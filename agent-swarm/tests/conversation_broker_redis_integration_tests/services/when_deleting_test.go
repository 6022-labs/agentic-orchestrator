package services_test

import (
	"context"

	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_redis/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	redis_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_redis_mocks/metrics_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenDeletingTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenDeletingBeforeEach(t *testing.T) *WhenDeletingTestingSuite {
	redisClient := container_builders.GetRedisContainerBuilder().Create()

	redisMetrics := redis_metrics_mocks.NewMockRedisMetricsInterface(gomock.NewController(t))
	redisMetrics.EXPECT().RecordOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	sut := services.NewRedisStorage(
		redisClient,
		redisMetrics,
	)

	return &WhenDeletingTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenDeletingAfterEach(suite *WhenDeletingTestingSuite) {
	container_builders.GetRedisContainerBuilder().DropConnection(suite.redisClient)
}

func TestWhenDeleting(t *testing.T) {
	t.Parallel()

	t.Run("Given key exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingBeforeEach(t)
			defer WhenDeletingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Store(context.Background(), key, "test-value", nil)
			assert.Nil(t, err)

			err = suite.sut.Delete(context.Background(), key)
			assert.Nil(t, err)
		})

		t.Run("Should delete key", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingBeforeEach(t)
			defer WhenDeletingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Store(context.Background(), key, "test-value", nil)
			assert.Nil(t, err)

			err = suite.sut.Delete(context.Background(), key)
			assert.Nil(t, err)

			value, err := suite.sut.Read(context.Background(), key)
			assert.NotNil(t, err)
			assert.Empty(t, value)
		})
	})

	t.Run("Given key does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingBeforeEach(t)
			defer WhenDeletingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Delete(context.Background(), key)
			assert.Nil(t, err)
		})
	})
}
