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
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"redis.host":     "localhost:6379",
			"redis.password": "secret",
			"redis.db":       int64(2),
		}, "."), nil)

		settings := redis_settings.NewRedisSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "localhost:6379", settings.Host)
		assert.Equal(t, "secret", settings.Password)
		assert.Equal(t, int64(2), settings.DB)
	})
}
