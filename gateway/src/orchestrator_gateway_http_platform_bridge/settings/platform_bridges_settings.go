package settings

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const PlatformBridgesSettingsKey = "platform_bridges"

type PlatformBridgesSettings struct {
	SlackBridgeUrl string `mapstructure:"slack_bridge_url"`
}

func NewPlatformBridgesSettings(logger *zap.Logger, v *viper.Viper) *PlatformBridgesSettings {
	if val, ok := os.LookupEnv("PLATFORM_BRIDGES__SLACK_BRIDGE_URL"); ok {
		v.Set(PlatformBridgesSettingsKey+".slack_bridge_url", val)
	}

	var settings PlatformBridgesSettings
	if err := v.UnmarshalKey(PlatformBridgesSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal platform bridges settings", zap.Error(err))
	}

	if len(strings.TrimSpace(settings.SlackBridgeUrl)) == 0 {
		logger.Fatal("please set your platform_bridges.slack_bridge_url value in your configuration")
	}

	return &settings
}
