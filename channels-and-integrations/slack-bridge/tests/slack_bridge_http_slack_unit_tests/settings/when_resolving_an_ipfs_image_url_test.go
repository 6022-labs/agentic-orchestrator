package settings_test

import (
	"testing"

	ipfs_gateway_settings "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_slack/settings"
	"github.com/stretchr/testify/assert"
)

func TestWhenResolvingAnIpfsImageUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given a {cid} template", func(t *testing.T) {
		t.Parallel()

		settings := &ipfs_gateway_settings.IpfsGatewaySettings{
			GatewayUrlTemplate: "https://gateway.example.com/ipfs/{cid}",
		}

		t.Run("Should substitute the cid", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "https://gateway.example.com/ipfs/QmHash", settings.ImageUrl("QmHash"))
		})

		t.Run("Should strip an ipfs scheme", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "https://gateway.example.com/ipfs/QmHash", settings.ImageUrl("ipfs://QmHash"))
		})

		t.Run("Should keep a path suffix out of the substitution", func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				"https://gateway.example.com/ipfs/QmHash/image.png",
				settings.ImageUrl("ipfs://QmHash/image.png"),
			)
		})

		t.Run("Should pass an already displayable url through", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "https://cdn.example.com/a.png", settings.ImageUrl("https://cdn.example.com/a.png"))
		})
	})

	t.Run("Given a legacy %s template", func(t *testing.T) {
		t.Parallel()

		settings := &ipfs_gateway_settings.IpfsGatewaySettings{
			GatewayUrlTemplate: "https://example.com/ipfs/%s?name=a%2Fb",
		}

		t.Run("Should substitute only the first verb and leave escapes alone", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "https://example.com/ipfs/QmHash?name=a%2Fb", settings.ImageUrl("QmHash"))
		})
	})

	t.Run("Given a template without a placeholder", func(t *testing.T) {
		t.Parallel()

		settings := &ipfs_gateway_settings.IpfsGatewaySettings{
			GatewayUrlTemplate: "https://example.com/ipfs/",
		}

		t.Run("Should append the cid", func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "https://example.com/ipfs/QmHash", settings.ImageUrl("QmHash"))
		})
	})
}
