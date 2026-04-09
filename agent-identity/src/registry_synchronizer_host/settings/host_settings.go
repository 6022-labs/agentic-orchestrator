package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const HostSettingsKey = "host"

type HostSettings struct {
	AppName string `koanf:"app_name"`
}

func NewHostSettings(logger *zap.Logger, k *koanf.Koanf) *HostSettings {
	settings := HostSettings{}
	if err := k.Unmarshal(HostSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal host settings", zap.Error(err))
	}

	settings.AppName = strings.TrimSpace(settings.AppName)

	if len(settings.AppName) == 0 {
		logger.Fatal("please set your host.app_name value in your configuration")
	}

	return &settings
}
