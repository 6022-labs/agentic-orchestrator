package services

import (
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_redis/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type WhenStoringTestSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenStoringBeforeEach() *WhenStoringTestSuite {
	redisClient := container_builders.GetRedisConnection()

	sut := services.NewRedisStorage(
		redisClient,
	)

	return &WhenStoringTestSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenStoringAfterEach(suite *WhenStoringTestSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenStoring(t *testing.T) {
	t.Parallel()

	t.Run("Given no expiration", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store value", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringBeforeEach()
			defer WhenStoringAfterEach(suite)

			err := suite.sut.Store("test-key-1", "test-value", nil)
			assert.Nil(t, err)
		})
	})

	t.Run("Given expiration", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store value", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringBeforeEach()
			defer WhenStoringAfterEach(suite)

			duration := time.Second

			err := suite.sut.Store("test-key-2", "test-value", &duration)
			assert.Nil(t, err)

			secret, err := suite.sut.Read("test-key-2")
			assert.Nil(t, err)
			assert.Equal(t, "test-value", secret)
		})

		t.Run("Should not be able to read value after expiration", func(t *testing.T) {
			t.Parallel()

			suite := WhenStoringBeforeEach()
			defer WhenStoringAfterEach(suite)

			duration := time.Second

			err := suite.sut.Store("test-key-3", "test-value", &duration)
			assert.Nil(t, err)

			time.Sleep(duration * 2)

			secret, err := suite.sut.Read("test-key-3")
			assert.NotNil(t, err)
			assert.Empty(t, secret)
		})
	})
}
