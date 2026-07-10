package services

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_redis/services"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/tests/common_tests/container_builders"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type WhenReadingTestSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenReadingBeforeEach() *WhenReadingTestSuite {
	redisClient := container_builders.GetRedisConnection()

	sut := services.NewRedisStorage(
		redisClient,
	)

	return &WhenReadingTestSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenReadingAfterEach(suite *WhenReadingTestSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenReading(t *testing.T) {
	t.Parallel()

	t.Run("Given key exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return value", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingBeforeEach()
			defer WhenReadingAfterEach(suite)

			key := uuid.New().String()
			err := suite.sut.Store(key, "test-value", nil)
			assert.Nil(t, err)

			value, err := suite.sut.Read(key)
			assert.Nil(t, err)
			assert.Equal(t, "test-value", value)
		})
	})

	t.Run("Given key does not exist", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenReadingBeforeEach()
			defer WhenReadingAfterEach(suite)

			key := uuid.New().String()
			value, err := suite.sut.Read(key)
			assert.NotNil(t, err)
			assert.Empty(t, value)
		})
	})
}
