package services

import (
	"context"

	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_redis/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	redis_metrics_mocks "github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_redis_mocks/metrics_mocks/interfaces_mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenStoringTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenStoringBeforeEach(t *testing.T) *WhenStoringTestingSuite {
	redisClient := container_builders.GetRedisContainerBuilder().Create()

	redisMetrics := redis_metrics_mocks.NewMockRedisMetricsInterface(gomock.NewController(t))
	redisMetrics.EXPECT().RecordOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	sut := services.NewRedisStorage(
		redisClient,
		redisMetrics,
	)

	return &WhenStoringTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenStoringAfterEach(suite *WhenStoringTestingSuite) {
	container_builders.GetRedisContainerBuilder().DropConnection(suite.redisClient)
}

func TestWhenStoring(t *testing.T) {
	t.Parallel()

	t.Run("Given no expiration", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store value", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringBeforeEach(t)
			defer WhenStoringAfterEach(suite)

			err := suite.sut.Store(context.Background(), "test-key-1", "test-value", nil)
			assert.Nil(t, err)
		})
	})

	t.Run("Given expiration", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store value", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringBeforeEach(t)
			defer WhenStoringAfterEach(suite)

			duration := time.Second

			err := suite.sut.Store(context.Background(), "test-key-2", "test-value", &duration)
			assert.Nil(t, err)

			secret, err := suite.sut.Read(context.Background(), "test-key-2")
			assert.Nil(t, err)
			assert.Equal(t, "test-value", secret)
		})

		t.Run("Should not be able to read value after expiration", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringBeforeEach(t)
			defer WhenStoringAfterEach(suite)

			duration := time.Second

			err := suite.sut.Store(context.Background(), "test-key-3", "test-value", &duration)
			assert.Nil(t, err)

			time.Sleep(duration * 2)

			secret, err := suite.sut.Read(context.Background(), "test-key-3")
			assert.NotNil(t, err)
			assert.Empty(t, secret)
		})
	})
}
