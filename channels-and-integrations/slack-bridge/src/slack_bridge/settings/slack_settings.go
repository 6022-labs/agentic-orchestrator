package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const SlackSettingsKey = "slack"

type SlackSettings struct {
	AppId                 string `koanf:"app_id"`
	ApiBaseURL            string `koanf:"api_base_url"`
	ClientSecret          string `koanf:"client_secret"`
	BotToken              string `koanf:"bot_token"`
	SigningSecret         string `koanf:"signing_secret"`
	AcceptedChannelPrefix string `koanf:"accepted_channel_prefix"`
}

func NewSlackSettings(logger *zap.Logger, k *koanf.Koanf) *SlackSettings {
	settings := SlackSettings{
		ApiBaseURL:            "https://slack.com",
		AcceptedChannelPrefix: "ai",
	}

	if err := k.Unmarshal(SlackSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal slack settings", zap.Error(err))
	}

	settings.AppId = strings.TrimSpace(settings.AppId)
	settings.ApiBaseURL = strings.TrimSpace(settings.ApiBaseURL)
	settings.ClientSecret = strings.TrimSpace(settings.ClientSecret)
	settings.BotToken = strings.TrimSpace(settings.BotToken)
	settings.SigningSecret = strings.TrimSpace(settings.SigningSecret)
	settings.AcceptedChannelPrefix = strings.TrimSpace(settings.AcceptedChannelPrefix)

	if len(settings.ApiBaseURL) == 0 {
		settings.ApiBaseURL = "https://slack.com"
	}

	if len(settings.AcceptedChannelPrefix) == 0 {
		settings.AcceptedChannelPrefix = "ai"
	}

	if len(strings.TrimSpace(settings.AppId)) == 0 {
		logger.Fatal("please set your slack.app_id value in your configuration")
	}

	if len(strings.TrimSpace(settings.ClientSecret)) == 0 {
		logger.Fatal("please set your slack.client_secret value in your configuration")
	}

	if len(strings.TrimSpace(settings.SigningSecret)) == 0 {
		logger.Fatal("please set your slack.signing_secret value in your configuration")
	}

	if len(strings.TrimSpace(settings.BotToken)) == 0 {
		logger.Fatal("please set your slack.bot_token value in your configuration")
	}

	return &settings
}
