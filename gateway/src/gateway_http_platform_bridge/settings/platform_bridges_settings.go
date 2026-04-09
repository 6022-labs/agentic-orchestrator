package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const PlatformBridgesSettingsKey = "platform_bridges"

type PlatformBridgesSettings struct {
	SlackBridgeUrl string `koanf:"slack_bridge_url"`
}

func NewPlatformBridgesSettings(logger *zap.Logger, k *koanf.Koanf) *PlatformBridgesSettings {
	var settings PlatformBridgesSettings
	if err := k.Unmarshal(PlatformBridgesSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal platform bridges settings", zap.Error(err))
	}

	settings.SlackBridgeUrl = strings.TrimSpace(settings.SlackBridgeUrl)

	if len(settings.SlackBridgeUrl) == 0 {
		logger.Fatal("please set your platform_bridges.slack_bridge_url value in your configuration")
	}

	return &settings
}
