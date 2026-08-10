package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const RegistryAgentResolutionSettingsKey = "registry_agent_resolution"

type RegistryAgentResolutionSettings struct {
	CacheTtlSeconds uint64 `koanf:"cache_ttl_seconds"`
}

func NewRegistryAgentResolutionSettings(logger *zap.Logger, k *koanf.Koanf) *RegistryAgentResolutionSettings {
	settings := RegistryAgentResolutionSettings{
		CacheTtlSeconds: 120,
	}
	if err := k.Unmarshal(RegistryAgentResolutionSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal registry agent resolution settings", zap.Error(err))
	}

	if settings.CacheTtlSeconds == 0 {
		settings.CacheTtlSeconds = 120
	}

	return &settings
}
