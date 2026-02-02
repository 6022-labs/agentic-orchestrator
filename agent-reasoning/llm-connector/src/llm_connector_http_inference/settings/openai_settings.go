package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const OpenAISettingsKey = "openai"

type OpenAISettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewOpenAISettings(logger *zap.Logger, v *viper.Viper) *OpenAISettings {
	v.SetDefault(OpenAISettingsKey+".base_url", "https://api.openai.com")
	if val, ok := os.LookupEnv("OPENAI__BASE_URL"); ok {
		v.Set(OpenAISettingsKey+".base_url", val)
	}

	var settings OpenAISettings
	if err := v.UnmarshalKey(OpenAISettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal openai settings", zap.Error(err))
	}

	return &settings
}

func (s *OpenAISettings) GetBaseUrl() string {
	return s.BaseUrl
}

func (s *OpenAISettings) HasDirectAccessToMcp() bool {
	return true
}

func (s *OpenAISettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIOpenAIResponsesStandard
}
