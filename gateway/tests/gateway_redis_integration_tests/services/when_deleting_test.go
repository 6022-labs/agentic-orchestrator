package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_redis/services"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests/container_builders"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type WhenDeletingTestSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenDeletingBeforeEach() *WhenDeletingTestSuite {
	redisClient := container_builders.GetRedisConnection()

	sut := services.NewRedisStorage(
		redisClient,
	)

	return &WhenDeletingTestSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenDeletingAfterEach(suite *WhenDeletingTestSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenDeleting(t *testing.T) {
	t.Parallel()

	t.Run("Given key exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingBeforeEach()
			defer WhenDeletingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Store(key, "test-value", nil)
			assert.Nil(t, err)

			err = suite.sut.Delete(key)
			assert.Nil(t, err)
		})

		t.Run("Should delete key", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingBeforeEach()
			defer WhenDeletingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Store(key, "test-value", nil)
			assert.Nil(t, err)

			err = suite.sut.Delete(key)
			assert.Nil(t, err)

			value, err := suite.sut.Read(key)
			assert.NotNil(t, err)
			assert.Empty(t, value)
		})
	})

	t.Run("Given key does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenDeletingBeforeEach()
			defer WhenDeletingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Delete(key)
			assert.Nil(t, err)
		})
	})
}
