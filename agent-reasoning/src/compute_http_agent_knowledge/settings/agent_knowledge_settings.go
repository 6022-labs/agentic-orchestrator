package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentKnowledgeSettingsKey = "agent_knowledge"

type AgentKnowledgeSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewAgentKnowledgeSettings(logger *zap.Logger, k *koanf.Koanf) *AgentKnowledgeSettings {
	var settings AgentKnowledgeSettings
	if err := k.Unmarshal(AgentKnowledgeSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent knowledge settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

	if len(settings.BaseUrl) == 0 {
		logger.Fatal("please set your agent_knowledge.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")

	return &settings
}
