package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentReasoningAgentComputeSettingsKey = "agent_reasoning_agent_compute"

type AgentReasoningAgentComputeSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewAgentReasoningAgentComputeSettings(logger *zap.Logger, k *koanf.Koanf) *AgentReasoningAgentComputeSettings {
	var settings AgentReasoningAgentComputeSettings
	if err := k.Unmarshal(AgentReasoningAgentComputeSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent_reasoning_agent_compute settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)
	if len(settings.BaseUrl) == 0 {
		logger.Fatal("please set your agent_reasoning_agent_compute.base_url value in your configuration")
	}
	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")
	return &settings
}
