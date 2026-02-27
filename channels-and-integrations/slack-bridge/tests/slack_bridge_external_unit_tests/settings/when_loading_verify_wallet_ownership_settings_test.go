package settings_test

import (
	"testing"

	verify_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingVerifyWalletOwnershipSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"verify_wallet_ownership.orchestrator_front_base_url": "https://front.example.com/",
		}, "."), nil)

		settings := verify_settings.NewVerifyWalletOwnershipSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://front.example.com", settings.OrchestratorFrontBaseUrl)
	})
}
