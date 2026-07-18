package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingBlockchainSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured chain id", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"blockchain.chain_id": 1,
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := blockchain_settings.NewBlockchainSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, int64(1), settings.ChainId)
		})
	})
}
