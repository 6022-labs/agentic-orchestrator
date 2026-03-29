package settings_test

import (
	"testing"

	agent_identity_registry_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_identity_registry/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentIdentityRegistrySettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_identity_registry.base_url": "https://agent-identity-registry.example.com/",
		}, "."), nil)

		settings := agent_identity_registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://agent-identity-registry.example.com", settings.BaseUrl)
	})
}
