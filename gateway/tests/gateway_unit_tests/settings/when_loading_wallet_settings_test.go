package settings_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingWalletSettings(t *testing.T) {
	t.Parallel()

	expectedAddress := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	t.Run("Given a 0x-prefixed private key", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"wallet.private_key": "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		}, "."), nil)

		t.Run("Should load it and derive the address", func(t *testing.T) {
			t.Parallel()

			walletSettings := settings.NewWalletSettings(zap.NewNop(), k)

			assert.NotNil(t, walletSettings)
			assert.Equal(
				t,
				"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
				walletSettings.PrivateKey,
			)
			assert.Equal(t, expectedAddress, walletSettings.Address)
		})
	})

	t.Run("Given a bare private key", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"wallet.private_key": " ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 ",
		}, "."), nil)

		t.Run("Should trim it and derive the address", func(t *testing.T) {
			t.Parallel()

			walletSettings := settings.NewWalletSettings(zap.NewNop(), k)

			assert.Equal(
				t,
				"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
				walletSettings.PrivateKey,
			)
			assert.Equal(t, expectedAddress, walletSettings.Address)
		})
	})
}
