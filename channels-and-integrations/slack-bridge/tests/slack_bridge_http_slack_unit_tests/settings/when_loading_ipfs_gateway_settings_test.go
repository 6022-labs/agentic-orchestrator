package settings_test

import (
	"testing"

	ipfs_gateway_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingIpfsGatewaySettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"ipfs_gateway.url_template": "https://gateway.example.com/ipfs/%s",
		}, "."), nil)

		settings := ipfs_gateway_settings.NewIpfsGatewaySettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://gateway.example.com/ipfs/%s", settings.GatewayUrlTemplate)
	})
}
