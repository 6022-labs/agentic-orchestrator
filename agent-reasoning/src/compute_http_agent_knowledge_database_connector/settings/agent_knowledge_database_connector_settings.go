package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentKnowledgeDatabaseConnectorSettingsKey = "agent_knowledge_database_connector"

type AgentKnowledgeDatabaseConnectorSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewAgentKnowledgeDatabaseConnectorSettings(logger *zap.Logger, k *koanf.Koanf) *AgentKnowledgeDatabaseConnectorSettings {
	var settings AgentKnowledgeDatabaseConnectorSettings
	if err := k.Unmarshal(AgentKnowledgeDatabaseConnectorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent knowledge database connector settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

	if len(settings.BaseUrl) == 0 {
		logger.Fatal("please set your agent_knowledge_database_connector.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")

	return &settings
}
