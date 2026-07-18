package settings_test

import (
	"testing"

	platform_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_platform_bridge/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingPlatformBridgesSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured slack bridge url", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"platform_bridges.slack_bridge_url": "https://slack.bridge",
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://slack.bridge", settings.SlackBridgeUrl)
		})
	})

	t.Run("Given a slack bridge url with surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"platform_bridges.slack_bridge_url": "  https://slack.bridge  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

			assert.Equal(t, "https://slack.bridge", settings.SlackBridgeUrl)
		})
	})
}
