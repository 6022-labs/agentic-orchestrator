package settings_test

import (
	"testing"

	ipfs_gateway_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingIpfsGatewaySettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("ipfs_gateway", map[string]any{
			"url_template": "https://gateway.example.com/ipfs/%s",
		})

		settings := ipfs_gateway_settings.NewIpfsGatewaySettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://gateway.example.com/ipfs/%s", settings.GatewayUrlTemplate)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("IPFS_GATEWAY__URL_TEMPLATE", "https://gateway.env/ipfs/%s")

		v := viper.New()
		settings := ipfs_gateway_settings.NewIpfsGatewaySettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://gateway.env/ipfs/%s", settings.GatewayUrlTemplate)
	})
}
