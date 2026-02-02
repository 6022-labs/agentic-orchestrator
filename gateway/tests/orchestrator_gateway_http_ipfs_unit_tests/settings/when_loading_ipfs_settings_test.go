package settings_test

import (
	"testing"

	ipfs_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_ipfs/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingIpfsSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("ipfs", map[string]any{
			"node_http_url": "https://ipfs.example.com",
		})

		settings := ipfs_settings.NewIPFSSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://ipfs.example.com", settings.IpfsNodeHttpUrl)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("IPFS__NODE_HTTP_URL", "https://ipfs.env")

		v := viper.New()
		settings := ipfs_settings.NewIPFSSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://ipfs.env", settings.IpfsNodeHttpUrl)
	})
}
