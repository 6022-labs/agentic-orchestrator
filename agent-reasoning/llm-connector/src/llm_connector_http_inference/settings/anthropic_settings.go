package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type AnthropicSettings struct {
	BaseUrl string
}

func NewAnthropicSettings(logger *zap.Logger) *AnthropicSettings {
	baseUrl := os.Getenv("ANTHROPIC_BASE_URL")
	if len(strings.TrimSpace(baseUrl)) == 0 {
		baseUrl = "https://api.anthropic.com"

		logger.Info("ANTHROPIC_BASE_URL is not set, using default base URL",
			zap.String("defaultBaseUrl", baseUrl),
		)
	}

	return &AnthropicSettings{
		BaseUrl: baseUrl,
	}
}

func (s *AnthropicSettings) GetBaseUrl() string {
	return s.BaseUrl
}

func (s *AnthropicSettings) HasDirectAccessToMcp() bool {
	return true
}

func (s *AnthropicSettings) PlatformAPI() PlatformAPIInterface {
	return PlatformAPIAnthropicMessagesStandard
}
