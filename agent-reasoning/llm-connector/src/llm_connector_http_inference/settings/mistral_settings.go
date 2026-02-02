package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const MistralSettingsKey = "mistral"

type MistralSettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewMistralSettings(logger *zap.Logger, v *viper.Viper) *MistralSettings {
	v.SetDefault(MistralSettingsKey+".base_url", "https://api.mistral.ai")
	if val, ok := os.LookupEnv("MISTRAL__BASE_URL"); ok {
		v.Set(MistralSettingsKey+".base_url", val)
	}

	var settings MistralSettings
	if err := v.UnmarshalKey(MistralSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal mistral settings", zap.Error(err))
	}

	return &settings
}

func (s *MistralSettings) GetBaseUrl() string {
	return s.BaseUrl
}

func (s *MistralSettings) HasDirectAccessToMcp() bool {
	return false
}

func (s *MistralSettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIOpenAICompletionsStandard
}
