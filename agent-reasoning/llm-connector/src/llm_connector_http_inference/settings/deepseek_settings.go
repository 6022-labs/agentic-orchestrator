package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const DeepSeekSettingsKey = "deepseek"

type DeepSeekSettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewDeepSeekSettings(logger *zap.Logger, v *viper.Viper) *DeepSeekSettings {
	v.SetDefault(DeepSeekSettingsKey+".base_url", "https://api.deepseek.com")
	if val, ok := os.LookupEnv("DEEPSEEK__BASE_URL"); ok {
		v.Set(DeepSeekSettingsKey+".base_url", val)
	}

	var settings DeepSeekSettings
	if err := v.UnmarshalKey(DeepSeekSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal deepseek settings", zap.Error(err))
	}

	return &settings
}

func (s *DeepSeekSettings) GetBaseUrl() string {
	return s.BaseUrl
}

func (s *DeepSeekSettings) HasDirectAccessToMcp() bool {
	return false
}

func (s *DeepSeekSettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIOpenAICompletionsStandard
}
