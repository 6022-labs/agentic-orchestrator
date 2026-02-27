package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingBlockchainSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"blockchain.chain_id": int64(1),
		}, "."), nil)

		settings := blockchain_settings.NewBlockchainSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, int64(1), settings.ChainId)
	})
}
