package settings_test

import (
	"testing"

	verify_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingVerifyWalletOwnershipSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("verify_wallet_ownership", map[string]any{
			"orchestrator_front_base_url": "https://front.example.com/",
		})

		settings := verify_settings.NewVerifyWalletOwnershipSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://front.example.com", settings.OrchestratorFrontBaseUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("VERIFY_WALLET_OWNERSHIP__ORCHESTRATOR_FRONT_BASE_URL", "https://front.env/")

		v := viper.New()
		settings := verify_settings.NewVerifyWalletOwnershipSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://front.env", settings.OrchestratorFrontBaseUrl)
	})
}
