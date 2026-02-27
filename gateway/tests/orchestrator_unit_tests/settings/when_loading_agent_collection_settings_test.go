package settings_test

import (
	"testing"

	agent_collection_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentCollectionSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_collection.signing_domain": "example.com",
		}, "."), nil)

		settings := agent_collection_settings.NewAgentCollectionSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "example.com", settings.SigningDomain)
	})
}
