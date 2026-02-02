package settings_test

import (
	"testing"

	agent_collection_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentCollectionSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("agent_collection", map[string]any{
			"signing_domain": "example.com",
		})

		settings := agent_collection_settings.NewAgentCollectionSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "example.com", settings.SigningDomain)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("AGENT_COLLECTION__SIGNING_DOMAIN", "env.example.com")

		v := viper.New()
		settings := agent_collection_settings.NewAgentCollectionSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "env.example.com", settings.SigningDomain)
	})
}
