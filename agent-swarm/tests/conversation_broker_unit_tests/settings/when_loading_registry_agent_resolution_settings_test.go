package settings_test

import (
	"testing"

	broker_settings "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRegistryAgentResolutionSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given configured values", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"registry_agent_resolution.cache_ttl_seconds": uint64(30),
		}, "."), nil)

		t.Run("Should load them", func(t *testing.T) {
			t.Parallel()

			settings := broker_settings.NewRegistryAgentResolutionSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(30), settings.CacheTtlSeconds)
		})
	})

	t.Run("Given the ttl is configured to zero", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"registry_agent_resolution.cache_ttl_seconds": uint64(0),
		}, "."), nil)

		t.Run("Should fall back to the default rather than cache forever", func(t *testing.T) {
			t.Parallel()

			settings := broker_settings.NewRegistryAgentResolutionSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(120), settings.CacheTtlSeconds)
		})
	})

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to defaults", func(t *testing.T) {
			t.Parallel()

			settings := broker_settings.NewRegistryAgentResolutionSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, uint64(120), settings.CacheTtlSeconds)
		})
	})
}
