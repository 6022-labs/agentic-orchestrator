package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const AnthropicSettingsKey = "anthropic"

type AnthropicSettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewAnthropicSettings(logger *zap.Logger, v *viper.Viper) *AnthropicSettings {
	v.SetDefault(AnthropicSettingsKey+".base_url", "https://api.anthropic.com")
	if val, ok := os.LookupEnv("ANTHROPIC__BASE_URL"); ok {
		v.Set(AnthropicSettingsKey+".base_url", val)
	}

	var settings AnthropicSettings
	if err := v.UnmarshalKey(AnthropicSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal anthropic settings", zap.Error(err))
	}

	return &settings
}

func (s *AnthropicSettings) GetBaseUrl() string {
	return s.BaseUrl
}

func (s *AnthropicSettings) HasDirectAccessToMcp() bool {
	return true
}

func (s *AnthropicSettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIAnthropicMessagesStandard
}
