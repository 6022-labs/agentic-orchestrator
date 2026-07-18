package settings_test

import (
	"testing"

	redis_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_redis/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRedisSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"redis.host":     "localhost:6379",
			"redis.password": "secret",
			"redis.db":       int64(2),
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := redis_settings.NewRedisSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "localhost:6379", settings.Host)
			assert.Equal(t, "secret", settings.Password)
			assert.Equal(t, int64(2), settings.DB)
		})
	})

	t.Run("Given values wrapped in whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"redis.host":     "  localhost:6379  ",
			"redis.password": "  secret  ",
		}, "."), nil)

		t.Run("Should trim them", func(t *testing.T) {
			t.Parallel()

			settings := redis_settings.NewRedisSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "localhost:6379", settings.Host)
			assert.Equal(t, "secret", settings.Password)
		})
	})

	t.Run("Given a negative db", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"redis.host": "localhost:6379",
			"redis.db":   int64(-3),
		}, "."), nil)

		t.Run("Should fall back to 0", func(t *testing.T) {
			t.Parallel()

			settings := redis_settings.NewRedisSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, int64(0), settings.DB)
		})
	})
}
