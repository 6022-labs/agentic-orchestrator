package settings_test

import (
	"testing"

	redis_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_redis/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRedisSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("redis", map[string]any{
			"host":     "localhost:6379",
			"password": "secret",
			"db":       int64(2),
		})

		settings := redis_settings.NewRedisSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "localhost:6379", settings.Host)
		assert.Equal(t, "secret", settings.Password)
		assert.Equal(t, int64(2), settings.DB)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("REDIS__HOST", "env:6379")
		t.Setenv("REDIS__PASSWORD", "env-secret")
		t.Setenv("REDIS__DB", "3")

		v := viper.New()
		settings := redis_settings.NewRedisSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "env:6379", settings.Host)
		assert.Equal(t, "env-secret", settings.Password)
		assert.Equal(t, int64(3), settings.DB)
	})
}
