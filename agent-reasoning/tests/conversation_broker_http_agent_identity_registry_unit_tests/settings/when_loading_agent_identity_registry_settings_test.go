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
	t.Parallel()

	t.Run("Given a configured base url", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_identity_registry.base_url": "https://agent-compute.example.com",
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://agent-compute.example.com", settings.BaseUrl)
		})
	})

	t.Run("Given a base url with a trailing slash", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_identity_registry.base_url": "https://agent-compute.example.com/",
		}, "."), nil)

		t.Run("Should trim the trailing slash", func(t *testing.T) {
			t.Parallel()

			settings := registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://agent-compute.example.com", settings.BaseUrl)
		})
	})
}
