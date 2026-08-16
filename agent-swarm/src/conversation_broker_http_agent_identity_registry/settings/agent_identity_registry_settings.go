package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentIdentityRegistrySettingsKey = "agent_identity_registry"

type AgentIdentityRegistrySettings struct {
	BaseUrl               string `koanf:"base_url"`
	RequestTimeoutSeconds uint64 `koanf:"request_timeout_seconds"`
}

func NewAgentIdentityRegistrySettings(logger *zap.Logger, k *koanf.Koanf) *AgentIdentityRegistrySettings {
	settings := AgentIdentityRegistrySettings{
		RequestTimeoutSeconds: 30,
	}
	if err := k.Unmarshal(AgentIdentityRegistrySettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent identity registry settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)
	if len(settings.BaseUrl) == 0 {
		logger.Fatal("please set your agent_identity_registry.base_url value in your configuration")
	}
	if settings.RequestTimeoutSeconds == 0 {
		settings.RequestTimeoutSeconds = 30
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")
	return &settings
}
