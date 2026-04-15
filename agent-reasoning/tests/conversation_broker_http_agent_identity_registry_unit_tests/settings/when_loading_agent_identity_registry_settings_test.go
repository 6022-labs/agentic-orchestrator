package settings_test

import (
	"testing"

	registry_settings "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_http_agent_identity_registry/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentIdentityRegistrySettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_identity_registry.base_url": "https://agent-compute.example.com",
		}, "."), nil)

		settings := registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-compute.example.com", settings.BaseUrl)
	})

	t.Run("Trims trailing slash from base_url", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_identity_registry.base_url": "https://agent-compute.example.com/",
		}, "."), nil)

		settings := registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-compute.example.com", settings.BaseUrl)
	})
}
