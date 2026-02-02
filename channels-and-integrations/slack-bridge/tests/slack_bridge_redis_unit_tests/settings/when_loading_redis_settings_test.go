package settings_test

import (
	"testing"

	redis_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_redis/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackBridgeRedisSettings(t *testing.T) {
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
		t.Setenv("REDIS__HOST", "localhost:6380")
		t.Setenv("REDIS__PASSWORD", "secret-env")
		t.Setenv("REDIS__DB", "3")

		v := viper.New()
		settings := redis_settings.NewRedisSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "localhost:6380", settings.Host)
		assert.Equal(t, "secret-env", settings.Password)
		assert.Equal(t, int64(3), settings.DB)
	})
}
