package settings_test

import (
	"testing"

	manager_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_blockchain/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingAgentCollectionsManagerSettings(t *testing.T) {
	address := "0x0000000000000000000000000000000000000002"

	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("agent_collections_manager", map[string]any{
			"smart_contract_address": address,
		})

		settings := manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, address, settings.SmartContractAddress.Hex())
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("AGENT_COLLECTIONS_MANAGER__SMART_CONTRACT_ADDRESS", address)

		v := viper.New()
		settings := manager_settings.NewAgentCollectionsManagerSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, address, settings.SmartContractAddress.Hex())
	})
}
