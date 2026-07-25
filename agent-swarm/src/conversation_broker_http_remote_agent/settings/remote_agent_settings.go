package settings

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const RemoteAgentSettingsKey = "remote_agent"

type RemoteAgentSettings struct {
	RequestTimeoutSeconds          uint64 `koanf:"request_timeout_seconds"`
	WellKnownRequestTimeoutSeconds uint64 `koanf:"well_known_request_timeout_seconds"`
}

func NewRemoteAgentSettings(logger *zap.Logger, k *koanf.Koanf) *RemoteAgentSettings {
	settings := RemoteAgentSettings{
		RequestTimeoutSeconds:          120,
		WellKnownRequestTimeoutSeconds: 10,
	}
	if err := k.Unmarshal(RemoteAgentSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent node settings", zap.Error(err))
	}

	if settings.RequestTimeoutSeconds == 0 {
		settings.RequestTimeoutSeconds = 120
	}
	if settings.WellKnownRequestTimeoutSeconds == 0 {
		settings.WellKnownRequestTimeoutSeconds = 10
	}

	return &settings
}
