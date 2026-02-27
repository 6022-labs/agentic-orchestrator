package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentCollectionSettingsKey = "agent_collection"

type AgentCollectionSettings struct {
	SigningDomain string `koanf:"signing_domain"`
}

func NewAgentCollectionSettings(logger *zap.Logger, k *koanf.Koanf) *AgentCollectionSettings {
	var settings AgentCollectionSettings
	if err := k.Unmarshal(AgentCollectionSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent collection settings", zap.Error(err))
	}

	settings.SigningDomain = strings.TrimSpace(settings.SigningDomain)

	if len(settings.SigningDomain) == 0 {
		logger.Fatal("please set your agent_collection.signing_domain value in your configuration")
	}

	return &settings
}
