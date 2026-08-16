package settings_test

import (
	"testing"

	platform_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_channels_and_integrations_bridge/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingPlatformBridgesSettings(t *testing.T) {
	t.Parallel()

	t.Run("Given a configured slack bridge", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"platform_bridges.slack_bridge.base_url": "https://slack.bridge",
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

			assert.NotNil(t, settings)
			assert.Equal(t, "https://slack.bridge", settings.SlackBridge.BaseUrl)
		})
	})

	t.Run("Given a slack bridge base url with a trailing slash", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"platform_bridges.slack_bridge.base_url": "https://slack.bridge/",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

			assert.Equal(t, "https://slack.bridge", settings.SlackBridge.BaseUrl)
		})
	})

	t.Run("Given a slack bridge base url with surrounding whitespace", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"platform_bridges.slack_bridge.base_url": "  https://slack.bridge  ",
		}, "."), nil)

		t.Run("Should trim it", func(t *testing.T) {
			t.Parallel()

			settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

			assert.Equal(t, "https://slack.bridge", settings.SlackBridge.BaseUrl)
		})
	})

	t.Run("Given no request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"platform_bridges.slack_bridge.base_url": "https://slack.bridge",
		}, "."), nil)

		t.Run("Should default to 30 seconds", func(t *testing.T) {
			t.Parallel()

			settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(30), settings.SlackBridge.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a zero request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"platform_bridges.slack_bridge.base_url":                "https://slack.bridge",
			"platform_bridges.slack_bridge.request_timeout_seconds": 0,
		}, "."), nil)

		t.Run("Should fall back to 30 seconds", func(t *testing.T) {
			t.Parallel()

			settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(30), settings.SlackBridge.RequestTimeoutSeconds)
		})
	})

	t.Run("Given a configured request timeout", func(t *testing.T) {
		t.Parallel()

		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]any{
			"platform_bridges.slack_bridge.base_url":                "https://slack.bridge",
			"platform_bridges.slack_bridge.request_timeout_seconds": 5,
		}, "."), nil)

		t.Run("Should load it", func(t *testing.T) {
			t.Parallel()

			settings := platform_settings.NewPlatformBridgesSettings(zap.NewNop(), k)

			assert.Equal(t, uint64(5), settings.SlackBridge.RequestTimeoutSeconds)
		})
	})
}
