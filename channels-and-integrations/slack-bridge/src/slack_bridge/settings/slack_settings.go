package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const SlackSettingsKey = "slack"

type SlackSettings struct {
	AppId                 string `mapstructure:"app_id"`
	ApiBaseURL            string `mapstructure:"api_base_url"`
	ClientSecret          string `mapstructure:"client_secret"`
	BotToken              string `mapstructure:"bot_token"`
	SigningSecret         string `mapstructure:"signing_secret"`
	AcceptedChannelPrefix string `mapstructure:"accepted_channel_prefix"`
}

func NewSlackSettings(logger *zap.Logger, v *viper.Viper) *SlackSettings {
	v.SetDefault(SlackSettingsKey+".api_base_url", "https://slack.com")
	v.SetDefault(SlackSettingsKey+".accepted_channel_prefix", "ai")

	if val, ok := os.LookupEnv("SLACK__APP_ID"); ok {
		v.Set(SlackSettingsKey+".app_id", val)
	}
	if val, ok := os.LookupEnv("SLACK__API_BASE_URL"); ok {
		v.Set(SlackSettingsKey+".api_base_url", val)
	}
	if val, ok := os.LookupEnv("SLACK__CLIENT_SECRET"); ok {
		v.Set(SlackSettingsKey+".client_secret", val)
	}
	if val, ok := os.LookupEnv("SLACK__SIGNING_SECRET"); ok {
		v.Set(SlackSettingsKey+".signing_secret", val)
	}
	if val, ok := os.LookupEnv("SLACK__BOT_TOKEN"); ok {
		v.Set(SlackSettingsKey+".bot_token", val)
	}
	if val, ok := os.LookupEnv("SLACK__ACCEPTED_CHANNEL_PREFIX"); ok {
		v.Set(SlackSettingsKey+".accepted_channel_prefix", val)
	}

	var settings SlackSettings
	if err := v.UnmarshalKey(SlackSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal slack settings", zap.Error(err))
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
