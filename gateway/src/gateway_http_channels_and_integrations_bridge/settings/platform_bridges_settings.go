package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const PlatformBridgesSettingsKey = "platform_bridges"

// PlatformBridgeSettings is one bridge's sub-configuration; each supported platform gets its own block.
type PlatformBridgeSettings struct {
	BaseUrl               string `koanf:"base_url"`
	RequestTimeoutSeconds uint64 `koanf:"request_timeout_seconds"`
}

type PlatformBridgesSettings struct {
	SlackBridge PlatformBridgeSettings `koanf:"slack_bridge"`
}

func NewPlatformBridgesSettings(logger *zap.Logger, k *koanf.Koanf) *PlatformBridgesSettings {
	settings := PlatformBridgesSettings{
		SlackBridge: PlatformBridgeSettings{
			RequestTimeoutSeconds: 30,
		},
	}
	if err := k.Unmarshal(PlatformBridgesSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal platform bridges settings", zap.Error(err))
	}

	settings.SlackBridge.BaseUrl = strings.TrimSpace(settings.SlackBridge.BaseUrl)

	if len(settings.SlackBridge.BaseUrl) == 0 {
		logger.Fatal("please set your platform_bridges.slack_bridge.base_url value in your configuration")
	}
	if settings.SlackBridge.RequestTimeoutSeconds == 0 {
		settings.SlackBridge.RequestTimeoutSeconds = 30
	}

	settings.SlackBridge.BaseUrl = strings.TrimRight(settings.SlackBridge.BaseUrl, "/")
	return &settings
}
