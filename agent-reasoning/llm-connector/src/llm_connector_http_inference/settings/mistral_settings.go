package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type MistralSettings struct {
	BaseUrl string
}

func NewMistralSettings(logger *zap.Logger) *MistralSettings {
	baseUrl := os.Getenv("MISTRAL_BASE_URL")
	if len(strings.TrimSpace(baseUrl)) == 0 {
		baseUrl = "https://api.mistral.ai"

		logger.Info("MISTRAL_BASE_URL is not set, using default base URL",
			zap.String("defaultBaseUrl", baseUrl),
		)
	}

	return &MistralSettings{
		BaseUrl: baseUrl,
	}
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
