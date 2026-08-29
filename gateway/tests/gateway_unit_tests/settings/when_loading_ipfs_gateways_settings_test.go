package settings_test

import (
	"testing"

	gateway_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingIpfsGatewaysSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given no configured templates", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")

		t.Run("Should fall back to the built-in gateways", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewIpfsGatewaysSettings(zap.NewNop(), k)

			assert.NotNil(t, result)
			assert.NotEmpty(t, result.GatewayUrlTemplates)
			assert.Contains(t, result.GatewayUrlTemplates, "https://w3s.link/ipfs/{cid}")
		})
	})

	t.Run("Given templates configured as a json array", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"ipfs.gateway_url_templates": []string{
				"https://{cid}.ipfs.example.com",
				"https://example.com/ipfs/{cid}",
			},
		}, "."), nil)

		t.Run("Should replace the built-in list in order", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewIpfsGatewaysSettings(zap.NewNop(), k)

			assert.Equal(t, []string{
				"https://{cid}.ipfs.example.com",
				"https://example.com/ipfs/{cid}",
			}, result.GatewayUrlTemplates)
		})
	})

	t.Run("Given templates configured as a comma separated string", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"ipfs.gateway_url_templates": " https://{cid}.ipfs.example.com , , https://example.com/ipfs/{cid} ",
		}, "."), nil)

		t.Run("Should split, trim and drop empty entries", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewIpfsGatewaysSettings(zap.NewNop(), k)

			assert.Equal(t, []string{
				"https://{cid}.ipfs.example.com",
				"https://example.com/ipfs/{cid}",
			}, result.GatewayUrlTemplates)
		})
	})

	t.Run("Given a blank override", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"ipfs.gateway_url_templates": "  ,  ",
		}, "."), nil)

		t.Run("Should fall back to the built-in gateways", func(t *testing.T) {
			t.Parallel()

			result := gateway_settings.NewIpfsGatewaysSettings(zap.NewNop(), k)

			assert.Contains(t, result.GatewayUrlTemplates, "https://w3s.link/ipfs/{cid}")
		})
	})
}
