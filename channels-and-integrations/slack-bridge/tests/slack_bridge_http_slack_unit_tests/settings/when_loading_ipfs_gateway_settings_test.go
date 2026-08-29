package settings_test

import (
	"testing"

	ipfs_gateway_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingIpfsGatewaySettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured template", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"ipfs_gateway.url_template": "https://gateway.example.com/ipfs/{cid}",
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := ipfs_gateway_settings.NewIpfsGatewaySettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://gateway.example.com/ipfs/{cid}", settings.GatewayUrlTemplate)
		})
	})

	t.Run("Given no configured template", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should default to a path gateway", func(t *testing.T) {
			t.Parallel()

			settings := ipfs_gateway_settings.NewIpfsGatewaySettings(zap.NewNop(), k)

			assert.Equal(t, "https://w3s.link/ipfs/{cid}", settings.GatewayUrlTemplate)
		})
	})
}
