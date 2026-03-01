package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const XSettingsKey = "x"

type XSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewXSettings(logger *zap.Logger, k *koanf.Koanf) *XSettings {
	settings := XSettings{BaseUrl: "https://api.x.ai"}
	if err := k.Unmarshal(XSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal x settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

	return &settings
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
