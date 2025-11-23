package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type OpenAISettings struct {
	BaseUrl string
}

func NewOpenAISettings(logger *zap.Logger) *OpenAISettings {
	baseUrl := os.Getenv("OPENAI_BASE_URL")
	if len(strings.TrimSpace(baseUrl)) == 0 {
		baseUrl = "https://api.openai.com"

		logger.Info("OPENAI_BASE_URL is not set, using default base URL",
			zap.String("defaultBaseUrl", baseUrl),
		)
	}

	return &OpenAISettings{
		BaseUrl: baseUrl,
	}
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
