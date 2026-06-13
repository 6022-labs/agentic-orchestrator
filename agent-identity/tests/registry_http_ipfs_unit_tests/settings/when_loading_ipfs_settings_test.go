package settings_test

import (
	"testing"

	ipfs_settings "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_http_ipfs/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingIPFSSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured node http url", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"ipfs.node_http_url": "https://ipfs.example.com",
		}, "."), nil)

		t.Run("Should load the url and leave the bearer token empty", func(t *testing.T) {
			t.Parallel()

			result := ipfs_settings.NewIPFSSettings(zap.NewNop(), k)

			assert.NotNil(t, result)
			assert.Equal(t, "https://ipfs.example.com", result.NodeHttpUrl)
			assert.Empty(t, result.BearerToken)
		})
	})

	t.Run("Given a configured bearer token", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"ipfs.node_http_url": "https://ipfs.example.com",
			"ipfs.bearer_token":  "  secret-token123  ",
		}, "."), nil)

		t.Run("Should load and trim the bearer token", func(t *testing.T) {
			t.Parallel()

			result := ipfs_settings.NewIPFSSettings(zap.NewNop(), k)

			assert.Equal(t, "secret-token123", result.BearerToken)
		})
	})

	t.Run("Given a url padded with whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"ipfs.node_http_url": "  http://127.0.0.1:5001  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			result := ipfs_settings.NewIPFSSettings(zap.NewNop(), k)

			assert.Equal(t, "http://127.0.0.1:5001", result.NodeHttpUrl)
		})
	})
}
