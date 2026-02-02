package settings_test

import (
	"testing"

	slack_settings "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/settings"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWhenLoadingSlackSettings(t *testing.T) {
	t.Run("Given settings provided through viper", func(t *testing.T) {
		v := viper.New()
		v.Set("slack", map[string]any{
			"app_id":                  "app-id",
			"api_base_url":            "https://slack.example.com",
			"client_secret":           "client-secret",
			"bot_token":               "bot-token",
			"signing_secret":          "signing-secret",
			"accepted_channel_prefix": "ai",
		})

		settings := slack_settings.NewSlackSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "app-id", settings.AppId)
		assert.Equal(t, "https://slack.example.com", settings.ApiBaseURL)
		assert.Equal(t, "client-secret", settings.ClientSecret)
		assert.Equal(t, "bot-token", settings.BotToken)
		assert.Equal(t, "signing-secret", settings.SigningSecret)
		assert.Equal(t, "ai", settings.AcceptedChannelPrefix)
	})

	t.Run("Given settings provided through env", func(t *testing.T) {
		t.Setenv("SLACK__APP_ID", "app-id-env")
		t.Setenv("SLACK__API_BASE_URL", "https://slack.env")
		t.Setenv("SLACK__CLIENT_SECRET", "client-secret-env")
		t.Setenv("SLACK__BOT_TOKEN", "bot-token-env")
		t.Setenv("SLACK__SIGNING_SECRET", "signing-secret-env")
		t.Setenv("SLACK__ACCEPTED_CHANNEL_PREFIX", "ai")

		v := viper.New()
		settings := slack_settings.NewSlackSettings(zap.NewNop(), v)

		assert.NotNil(t, settings)
		assert.Equal(t, "app-id-env", settings.AppId)
		assert.Equal(t, "https://slack.env", settings.ApiBaseURL)
		assert.Equal(t, "client-secret-env", settings.ClientSecret)
		assert.Equal(t, "bot-token-env", settings.BotToken)
		assert.Equal(t, "signing-secret-env", settings.SigningSecret)
		assert.Equal(t, "ai", settings.AcceptedChannelPrefix)
	})
}
