package services_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_redis/services"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/tests/common_tests/container_builders"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type WhenReadingBoolTestingSuite struct {
	redisClient *redis.Client
	sut         *services.RedisStorage
}

func WhenReadingBoolBeforeEach() *WhenReadingBoolTestingSuite {
	redisClient := container_builders.GetRedisConnection()

	sut := services.NewRedisStorage(
		redisClient,
	)

	return &WhenReadingBoolTestingSuite{
		redisClient: redisClient,
		sut:         sut,
	}
}

func WhenReadingBoolAfterEach(suite *WhenReadingBoolTestingSuite) {
	container_builders.CloseRedisConnection(suite.redisClient)
}

func TestWhenReadingBool(t *testing.T) {
	t.Parallel()

	t.Run("Given a key storing 'true' as string, When reading it as bool, Then it returns true", func(t *testing.T) {
		key := fmt.Sprintf("test_key_true_%s_%d", t.Name(), time.Now().UnixNano())

		suite := WhenReadingBoolBeforeEach()
		defer WhenReadingBoolAfterEach(suite)

		err := suite.redisClient.Set(context.Background(), key, "true", 0).Err()
		if err != nil {
			t.Fatalf("Failed to set key in Redis: %v", err)
		}

		result, err := suite.sut.ReadBool(key)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, *result)
	})

	t.Run("Given a key storing 'false' as string, When reading it as bool, Then it returns false", func(t *testing.T) {
		key := fmt.Sprintf("test_key_false_%s_%d", t.Name(), time.Now().UnixNano())

		suite := WhenReadingBoolBeforeEach()
		defer WhenReadingBoolAfterEach(suite)

		err := suite.redisClient.Set(context.Background(), key, "false", 0).Err()
		if err != nil {
			t.Fatalf("Failed to set key in Redis: %v", err)
		}

		result, err := suite.sut.ReadBool(key)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, *result)
	})

	t.Run("Given a key storing an invalid bool string, When reading it as bool, Then it returns an error", func(t *testing.T) {
		key := fmt.Sprintf("test_key_invalid_%s_%d", t.Name(), time.Now().UnixNano())

		suite := WhenReadingBoolBeforeEach()
		defer WhenReadingBoolAfterEach(suite)

		err := suite.redisClient.Set(context.Background(), key, "not_a_bool", 0).Err()
		if err != nil {
			t.Fatalf("Failed to set key in Redis: %v", err)
		}

		_, err = suite.sut.ReadBool(key)

		assert.Error(t, err)
	})

	t.Run("Given a non-existing key, When reading it as bool, Then it returns nothing", func(t *testing.T) {
		key := fmt.Sprintf("non_existing_%s_%d", t.Name(), time.Now().UnixNano())

		suite := WhenReadingBoolBeforeEach()
		defer WhenReadingBoolAfterEach(suite)

		result, err := suite.sut.ReadBool(key)

		assert.Nil(t, err)
		assert.Nil(t, result)
	})
}
