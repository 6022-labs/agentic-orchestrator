package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const LibertAISettingsKey = "libertai"

type LibertAISettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewLibertAISettings(logger *zap.Logger, v *viper.Viper) *LibertAISettings {
	v.SetDefault(LibertAISettingsKey+".base_url", "https://api.libertai.io")
	if val, ok := os.LookupEnv("LIBERTAI__BASE_URL"); ok {
		v.Set(LibertAISettingsKey+".base_url", val)
	}

	var settings LibertAISettings
	if err := v.UnmarshalKey(LibertAISettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal libertai settings", zap.Error(err))
	}

	return &settings
}

func (s *LibertAISettings) GetBaseUrl() string {
	return s.BaseUrl
}

func (s *LibertAISettings) HasDirectAccessToMcp() bool {
	return false
}

func (s *LibertAISettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIOpenAICompletionsStandard
}
