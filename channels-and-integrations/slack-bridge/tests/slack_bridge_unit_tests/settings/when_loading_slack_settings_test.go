package settings_test

import (
	"testing"

	slack_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackSettings(t *testing.T) {
	t.Run("Loads settings", func(t *testing.T) {
		k := koanf.New(".")
		_ = k.Load(confmap.Provider(map[string]interface{}{
			"slack.app_id":                  "app-id",
			"slack.api_base_url":            "https://slack.example.com",
			"slack.client_secret":           "client-secret",
			"slack.bot_token":               "bot-token",
			"slack.signing_secret":          "signing-secret",
			"slack.accepted_channel_prefix": "ai",
		}, "."), nil)

		settings := slack_settings.NewSlackSettings(zap.NewNop(), k)

		assert.NotNil(t, settings)
		assert.Equal(t, "app-id", settings.AppId)
		assert.Equal(t, "https://slack.example.com", settings.ApiBaseURL)
		assert.Equal(t, "client-secret", settings.ClientSecret)
		assert.Equal(t, "bot-token", settings.BotToken)
		assert.Equal(t, "signing-secret", settings.SigningSecret)
		assert.Equal(t, "ai", settings.AcceptedChannelPrefix)
	})
}
