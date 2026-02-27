package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const OpenAISettingsKey = "openai"

type OpenAISettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewOpenAISettings(logger *zap.Logger, k *koanf.Koanf) *OpenAISettings {
	settings := OpenAISettings{BaseUrl: "https://api.openai.com"}
	if err := k.Unmarshal(OpenAISettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal openai settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

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
