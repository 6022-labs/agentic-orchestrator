package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_redis/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type WhenReadingAndDeletingTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenReadingAndDeletingBeforeEach() *WhenReadingAndDeletingTestingSuite {
	redisClient := container_builders.GetRedisConnection()

	sut := services.NewRedisStorage(
		redisClient,
	)

	return &WhenReadingAndDeletingTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenReadingAndDeletingAfterEach(suite *WhenReadingAndDeletingTestingSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenReadingAndDeleting(t *testing.T) {
	t.Parallel()

	t.Run("Given key exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the value and remove the key", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingAndDeletingBeforeEach()
			defer WhenReadingAndDeletingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Store(key, "test-value", nil)
			assert.Nil(t, err)

			value, err := suite.sut.ReadAndDelete(key)
			assert.Nil(t, err)
			assert.Equal(t, "test-value", value)

			_, err = suite.sut.ReadAndDelete(key)
			assert.NotNil(t, err)
		})
	})

	t.Run("Given key does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingAndDeletingBeforeEach()
			defer WhenReadingAndDeletingAfterEach(suite)

			_, err := suite.sut.ReadAndDelete(uuid.New().String())
			assert.NotNil(t, err)
		})
	})
}
