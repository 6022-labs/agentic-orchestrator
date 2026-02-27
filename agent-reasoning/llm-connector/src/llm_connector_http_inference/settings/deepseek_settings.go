package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const DeepSeekSettingsKey = "deepseek"

type DeepSeekSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewDeepSeekSettings(logger *zap.Logger, k *koanf.Koanf) *DeepSeekSettings {
	settings := DeepSeekSettings{BaseUrl: "https://api.deepseek.com"}
	if err := k.Unmarshal(DeepSeekSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal deepseek settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

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
