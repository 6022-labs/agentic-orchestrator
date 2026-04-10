package settings_test

import (
	"testing"

	rate_limit_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRateLimitSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"rate_limit.requests_per_second": 120,
		}, "."), nil)

		settings := rate_limit_settings.NewRateLimitSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 120, settings.RequestsPerSecond)
	})

	t.Run("Uses default when unset", func(t *testing.T) {
		k := koanf.New(".")

		settings := rate_limit_settings.NewRateLimitSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 50, settings.RequestsPerSecond)
	})

	t.Run("Falls back to default when non-positive", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"rate_limit.requests_per_second": 0,
		}, "."), nil)

		settings := rate_limit_settings.NewRateLimitSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, 50, settings.RequestsPerSecond)
	})
}
