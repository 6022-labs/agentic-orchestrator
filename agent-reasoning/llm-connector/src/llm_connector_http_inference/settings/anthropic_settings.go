package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AnthropicSettingsKey = "anthropic"

type AnthropicSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewAnthropicSettings(logger *zap.Logger, k *koanf.Koanf) *AnthropicSettings {
	settings := AnthropicSettings{BaseUrl: "https://api.anthropic.com"}
	if err := k.Unmarshal(AnthropicSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal anthropic settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

	return &settings
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
