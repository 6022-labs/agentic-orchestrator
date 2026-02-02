package settings_test

import (
	"testing"

	jwt_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingJwtSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("jwt", map[string]any{
			"hash":        "hash-value",
			"living_time": 3600,
		})

		settings := jwt_settings.NewJwtSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "hash-value", settings.Hash)
		assert.Equal(t, 3600, settings.LivingTime)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("JWT__HASH", "env-hash")
		t.Setenv("JWT__LIVING_TIME", "7200")

		v := viper.New()
		settings := jwt_settings.NewJwtSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "env-hash", settings.Hash)
		assert.Equal(t, 7200, settings.LivingTime)
	})
}
