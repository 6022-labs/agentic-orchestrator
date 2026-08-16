package settings_test

import (
	"testing"

	agent_identity_registry_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_http_agent_identity_registry/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentIdentityRegistrySettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured base url with a trailing slash", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_identity_registry.base_url": "https://agent-identity-registry.example.com/",
		}, "."), nil)

		t.Run("Should load and trim it", func(t *testing.T) {
			t.Parallel()

			settings := agent_identity_registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://agent-identity-registry.example.com", settings.BaseUrl)
		})
	})

	t.Run("Given a base url with surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_identity_registry.base_url": "  https://agent-identity-registry.example.com/  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			settings := agent_identity_registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://agent-identity-registry.example.com", settings.BaseUrl)
		})
	})

	t.Run("Given no request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_identity_registry.base_url": "https://agent-identity-registry.example.com",
		}, "."), nil)

		t.Run("Should default to 30 seconds", func(t *testing.T) {
			t.Parallel()

			settings := agent_identity_registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

			assert.Equal(t, uint64(30), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a zero request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_identity_registry.base_url":                "https://agent-identity-registry.example.com",
			"agent_identity_registry.request_timeout_seconds": 0,
		}, "."), nil)

		t.Run("Should fall back to 30 seconds", func(t *testing.T) {
			t.Parallel()

			settings := agent_identity_registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

			assert.Equal(t, uint64(30), settings.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a configured request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"agent_identity_registry.base_url":                "https://agent-identity-registry.example.com",
			"agent_identity_registry.request_timeout_seconds": 5,
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := agent_identity_registry_settings.NewAgentIdentityRegistrySettings(zap.NewNop(), k)

			assert.Equal(t, uint64(5), settings.RequestTimeoutSeconds)
		})
	})
}
