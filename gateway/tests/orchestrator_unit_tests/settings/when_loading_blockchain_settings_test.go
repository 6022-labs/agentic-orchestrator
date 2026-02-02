package settings_test

import (
	"testing"

	blockchain_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingBlockchainSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("blockchain", map[string]any{
			"chain_id": int64(1),
		})

		settings := blockchain_settings.NewBlockchainSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, int64(1), settings.ChainId)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("BLOCKCHAIN__CHAIN_ID", "10")

		v := viper.New()
		settings := blockchain_settings.NewBlockchainSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, int64(10), settings.ChainId)
	})
}
