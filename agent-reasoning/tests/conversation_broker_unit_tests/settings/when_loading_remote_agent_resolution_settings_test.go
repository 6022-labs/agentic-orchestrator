package settings_test

import (
	"testing"

	broker_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingRemoteAgentResolutionSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"remote_agent_resolution.cache_ttl_seconds": uint64(60),
		}, "."), nil)

		settings := broker_settings.NewRemoteAgentResolutionSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(60), settings.CacheTtlSeconds)
	})

	t.Run("Defaults missing values", func(t *testing.T) {
		k := koanf.New(".")

		settings := broker_settings.NewRemoteAgentResolutionSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, uint64(300), settings.CacheTtlSeconds)
	})
}
