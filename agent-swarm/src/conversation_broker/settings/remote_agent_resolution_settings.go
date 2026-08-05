package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const RemoteAgentResolutionSettingsKey = "remote_agent_resolution"

type RemoteAgentResolutionSettings struct {
	CacheTtlSeconds uint64 `koanf:"cache_ttl_seconds"`
}

func NewRemoteAgentResolutionSettings(logger *zap.Logger, k *koanf.Koanf) *RemoteAgentResolutionSettings {
	settings := RemoteAgentResolutionSettings{
		CacheTtlSeconds: 300,
	}
	if err := k.Unmarshal(RemoteAgentResolutionSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal remote agent resolution settings", zap.Error(err))
	}

	if settings.CacheTtlSeconds == 0 {
		settings.CacheTtlSeconds = 300
	}

	return &settings
}
