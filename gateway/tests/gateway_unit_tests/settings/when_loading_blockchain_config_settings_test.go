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
	t.Parallel()

	t.Run("Given explicit addresses for a built-in chain", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"blockchain.evm_chains.80002.collections_manager_address": "0x39DD75a025145770a8b8278Cc8BA705116443Ea7",
			"blockchain.evm_chains.80002.collection_creator_address":  "0x99065eCD99aE9fBA4dC467C1588b9205Aa71Aa16",
		}, "."), nil)

		t.Run("Should keep them over the built-in deployment", func(t *testing.T) {
			t.Parallel()

			settings := gateway_settings.NewBlockchainConfigSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			if assert.Len(t, settings.EvmChains, 1) {
				assert.Equal(t, "0x39DD75a025145770a8b8278Cc8BA705116443Ea7", settings.EvmChains["80002"].CollectionsManagerAddress)
				assert.Equal(t, "0x99065eCD99aE9fBA4dC467C1588b9205Aa71Aa16", settings.EvmChains["80002"].CollectionCreatorAddress)
				assert.Equal(t, "https://amoy.polygonscan.com", settings.EvmChains["80002"].ScannerBaseUrl)
			}
		})
	})

	t.Run("Given a configured chain with missing address fields", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"blockchain.evm_chains.80002.scanner_base_url": "https://scanner.example.com",
		}, "."), nil)

		t.Run("Should fill them from the built-in deployment", func(t *testing.T) {
			t.Parallel()

			settings := gateway_settings.NewBlockchainConfigSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			if assert.Len(t, settings.EvmChains, 1) {
				assert.NotEmpty(t, settings.EvmChains["80002"].CollectionsManagerAddress)
				assert.NotEmpty(t, settings.EvmChains["80002"].CollectionCreatorAddress)
				assert.Equal(t, "https://scanner.example.com", settings.EvmChains["80002"].ScannerBaseUrl)
			}
		})
	})

	t.Run("Given no configuration", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to the built-in chains with a complete deployment", func(t *testing.T) {
			t.Parallel()

			settings := gateway_settings.NewBlockchainConfigSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			if assert.Len(t, settings.EvmChains, 2) {
				assert.NotEmpty(t, settings.EvmChains["80002"].CollectionsManagerAddress)
				assert.NotEmpty(t, settings.EvmChains["80002"].CollectionCreatorAddress)
				assert.Equal(t, "https://amoy.polygonscan.com", settings.EvmChains["80002"].ScannerBaseUrl)

				assert.NotEmpty(t, settings.EvmChains["137"].CollectionsManagerAddress)
				assert.NotEmpty(t, settings.EvmChains["137"].CollectionCreatorAddress)
				assert.Equal(t, "https://polygonscan.com", settings.EvmChains["137"].ScannerBaseUrl)
			}
		})
	})
}
