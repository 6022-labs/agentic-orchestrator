package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const HostSettingsKey = "host"

type HostSettings struct {
	UseApi         bool   `koanf:"use_api"`
	UseListeners   bool   `koanf:"use_listeners"`
	ApiPort        int    `koanf:"api_port"`
	AppName        string `koanf:"app_name"`
	BodyLimitBytes int    `koanf:"body_limit_bytes"`
}

func NewHostSettings(logger *zap.Logger, k *koanf.Koanf) *HostSettings {
	settings := HostSettings{
		UseApi:         true,
		UseListeners:   true,
		ApiPort:        3000,
		BodyLimitBytes: 4 * 1024 * 1024,
	}
	if err := k.Unmarshal(HostSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal host settings", zap.Error(err))
	}

	settings.AppName = strings.TrimSpace(settings.AppName)

	if settings.ApiPort <= 0 {
		logger.Warn("invalid host.api_port value, defaulting to 3000")
		settings.ApiPort = 3000
	}

	if len(settings.AppName) == 0 {
		logger.Fatal("please set your host.app_name value in your configuration")
	}

	if settings.BodyLimitBytes <= 0 {
		logger.Warn("invalid host.body_limit_bytes value, defaulting to 4MB")
		settings.BodyLimitBytes = 4 * 1024 * 1024
	}

	return &settings
}
