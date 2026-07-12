package settings_test

import (
	"testing"

	gateway_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingBlockchainConfigSettings(t *testing.T) {
	t.Run("Loads chains keyed by chain id", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"blockchain.evm_chains.80002.collections_manager_address": "0x39DD75a025145770a8b8278Cc8BA705116443Ea7",
			"blockchain.evm_chains.80002.collection_creator_address":  "0x99065eCD99aE9fBA4dC467C1588b9205Aa71Aa16",
		}, "."), nil)

		settings := gateway_settings.NewBlockchainConfigSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		if assert.Len(t, settings.EvmChains, 1) {
			assert.Equal(t, "0x39DD75a025145770a8b8278Cc8BA705116443Ea7", settings.EvmChains["80002"].CollectionsManagerAddress)
			assert.Equal(t, "0x99065eCD99aE9fBA4dC467C1588b9205Aa71Aa16", settings.EvmChains["80002"].CollectionCreatorAddress)
		}
	})

	t.Run("Defaults to an empty map", func(t *testing.T) {
		k := koanf.New(".")

		settings := gateway_settings.NewBlockchainConfigSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Empty(t, settings.EvmChains)
	})
}
