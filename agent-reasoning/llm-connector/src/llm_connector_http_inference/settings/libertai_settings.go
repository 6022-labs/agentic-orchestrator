package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const LibertAISettingsKey = "libertai"

type LibertAISettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewLibertAISettings(logger *zap.Logger, k *koanf.Koanf) *LibertAISettings {
	settings := LibertAISettings{BaseUrl: "https://api.libertai.io"}
	if err := k.Unmarshal(LibertAISettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal libertai settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

	return &settings
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
