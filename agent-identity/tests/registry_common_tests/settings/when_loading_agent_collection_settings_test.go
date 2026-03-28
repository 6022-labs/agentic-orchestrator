package settings_test

import (
	"testing"

	agent_collection_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentCollectionSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_collection.signing_domain": "AgentCollection",
		}, "."), nil)

		settings := agent_collection_settings.NewAgentCollectionSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "AgentCollection", settings.SigningDomain)
	})
}
