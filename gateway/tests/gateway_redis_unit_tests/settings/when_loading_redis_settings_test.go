package settings_test

import (
	"testing"

	redis_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_redis/settings"
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

	t.Run("Given a negative db index", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"redis.host": "localhost:6379",
			"redis.db":   int64(-1),
		}, "."), nil)

		t.Run("Should fall back to db 0", func(t *testing.T) {
			t.Parallel()

			settings := redis_settings.NewRedisSettings(zap.NewNop(), k)

			assert.Equal(t, int64(0), settings.DB)
		})
	})
}
