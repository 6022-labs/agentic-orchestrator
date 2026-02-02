package settings_test

import (
	"testing"

	orchestrator_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingOrchestratorSlackBridgeSettings(t *testing.T) {
	address := "0x0000000000000000000000000000000000000001"

	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("orchestrator_slack_bridge", map[string]any{
			"public_key": address,
		})

		settings := orchestrator_settings.NewOrchestratorSlackBridgeSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.True(t, common.HexToAddress(address) == settings.PublicKey)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("ORCHESTRATOR_SLACK_BRIDGE__PUBLIC_KEY", address)

		v := viper.New()
		settings := orchestrator_settings.NewOrchestratorSlackBridgeSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.True(t, common.HexToAddress(address) == settings.PublicKey)
	})
}
