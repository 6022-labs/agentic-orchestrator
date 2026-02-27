package settings_test

import (
	"testing"

	manager_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_blockchain/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentCollectionsManagerSettings(t *testing.T) {
	address := "0x0000000000000000000000000000000000000002"

	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"agent_collections_manager.smart_contract_address": address,
		}, "."), nil)

		settings := manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, address, settings.SmartContractAddress.Hex())
	})
}
