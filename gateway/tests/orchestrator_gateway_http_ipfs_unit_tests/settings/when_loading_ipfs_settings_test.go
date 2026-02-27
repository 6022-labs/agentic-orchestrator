package settings_test

import (
	"testing"

	ipfs_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator_gateway_http_ipfs/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingIpfsSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"ipfs.node_http_url": "https://ipfs.example.com",
		}, "."), nil)

		settings := ipfs_settings.NewIPFSSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "https://ipfs.example.com", settings.IpfsNodeHttpUrl)
	})
}
