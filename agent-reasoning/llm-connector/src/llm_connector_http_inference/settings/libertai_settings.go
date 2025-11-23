package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type LibertAISettings struct {
	BaseUrl string
}

func NewLibertAISettings(logger *zap.Logger) *LibertAISettings {
	baseUrl := os.Getenv("LIBERTAI_BASE_URL")
	if len(strings.TrimSpace(baseUrl)) == 0 {
		baseUrl = "https://api.libertai.io"

		logger.Info("LIBERTAI_BASE_URL is not set, using default base URL",
			zap.String("defaultBaseUrl", baseUrl),
		)
	}

	return &LibertAISettings{
		BaseUrl: baseUrl,
	}
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
