package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentDatabaseConnectorSettingsKey = "agent_database_connector"

type AgentDatabaseConnectorSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewAgentDatabaseConnectorSettings(logger *zap.Logger, k *koanf.Koanf) *AgentDatabaseConnectorSettings {
	var settings AgentDatabaseConnectorSettings
	if err := k.Unmarshal(AgentDatabaseConnectorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent database connector settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

	if len(settings.BaseUrl) == 0 {
		logger.Fatal("please set your agent_database_connector.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")

	return &settings
}
