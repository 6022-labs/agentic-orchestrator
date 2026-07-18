package settings_test

import (
	"testing"

	jwt_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingJwtSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"jwt.hash":        "test-signing-hash-at-least-32-characters",
			"jwt.living_time": 3600,
		}, "."), nil)

		settings := jwt_settings.NewJwtSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "test-signing-hash-at-least-32-characters", settings.Hash)
		assert.Equal(t, 3600, settings.LivingTime)
		assert.Equal(t, 300, settings.NonceTtl)
	})
}
