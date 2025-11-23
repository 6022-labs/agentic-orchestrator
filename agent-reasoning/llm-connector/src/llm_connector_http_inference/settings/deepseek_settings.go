package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type DeepSeekSettings struct {
	BaseUrl string
}

func NewDeepSeekSettings(logger *zap.Logger) *DeepSeekSettings {
	baseUrl := os.Getenv("DEEPSEEK_BASE_URL")
	if len(strings.TrimSpace(baseUrl)) == 0 {
		baseUrl = "https://api.deepseek.com"

		logger.Info("DEEPSEEK_BASE_URL is not set, using default base URL",
			zap.String("defaultBaseUrl", baseUrl),
		)
	}

	return &DeepSeekSettings{
		BaseUrl: baseUrl,
	}
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
