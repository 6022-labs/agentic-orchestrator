package settings_test

import (
	"testing"

	host_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_host/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingHostSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"host.api_port":         9201,
			"host.body_limit_bytes": 1024,
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := host_settings.NewHostSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, 9201, settings.ApiPort)
			assert.Equal(t, 1024, settings.BodyLimitBytes)
		})
	})

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to defaults", func(t *testing.T) {
			t.Parallel()

			settings := host_settings.NewHostSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, 3000, settings.ApiPort)
			assert.Equal(t, 4*1024*1024, settings.BodyLimitBytes)
		})
	})
}
