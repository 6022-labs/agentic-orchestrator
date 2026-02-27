package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentLLMConnectorSettingsKey = "agent_llm_connector"

type AgentLLMConnectorSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewAgentLLMConnectorSettings(logger *zap.Logger, k *koanf.Koanf) *AgentLLMConnectorSettings {
	var settings AgentLLMConnectorSettings
	if err := k.Unmarshal(AgentLLMConnectorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent llm connector settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)

	if len(settings.BaseUrl) == 0 {
		logger.Fatal("please set your agent_llm_connector.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")

	return &settings
}
