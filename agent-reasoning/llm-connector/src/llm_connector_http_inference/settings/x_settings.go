package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type XSettings struct {
	BaseUrl string
}

func NewXSettings(logger *zap.Logger) *XSettings {
	baseUrl := os.Getenv("X_BASE_URL")
	if len(strings.TrimSpace(baseUrl)) == 0 {
		baseUrl = "https://api.x.ai"

		logger.Info("X_BASE_URL is not set, using default base URL",
			zap.String("defaultBaseUrl", baseUrl),
		)
	}

	return &XSettings{
		BaseUrl: baseUrl,
	}
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
