package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const MistralSettingsKey = "mistral"

type MistralSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewMistralSettings(logger *zap.Logger, k *koanf.Koanf) *MistralSettings {
	settings := MistralSettings{BaseUrl: "https://api.mistral.ai"}
	if err := k.Unmarshal(MistralSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal mistral settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

	return &settings
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
