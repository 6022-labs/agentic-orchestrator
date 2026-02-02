package settings

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const XSettingsKey = "x"

type XSettings struct {
	BaseUrl string `mapstructure:"base_url"`
}

func NewXSettings(logger *zap.Logger, v *viper.Viper) *XSettings {
	v.SetDefault(XSettingsKey+".base_url", "https://api.x.ai")
	if val, ok := os.LookupEnv("X__BASE_URL"); ok {
		v.Set(XSettingsKey+".base_url", val)
	}

	var settings XSettings
	if err := v.UnmarshalKey(XSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal x settings", zap.Error(err))
	}

	return &settings
}

func (s *XSettings) GetBaseUrl() string {
	return s.BaseUrl
}

func (s *XSettings) HasDirectAccessToMcp() bool {
	return false
}

func (s *XSettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIOpenAICompletionsStandard
}
