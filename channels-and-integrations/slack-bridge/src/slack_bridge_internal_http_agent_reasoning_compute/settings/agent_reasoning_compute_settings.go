package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const AgentReasoningComputeSettingsKey = "agent_reasoning_compute"

type AgentReasoningComputeSettings struct {
	BaseUrl string `koanf:"base_url"`
}

func NewAgentReasoningComputeSettings(logger *zap.Logger, k *koanf.Koanf) *AgentReasoningComputeSettings {
	var settings AgentReasoningComputeSettings
	if err := k.Unmarshal(AgentReasoningComputeSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal agent reasoning compute settings", zap.Error(err))
	}

	settings.BaseUrl = strings.TrimSpace(settings.BaseUrl)
	if len(settings.BaseUrl) == 0 {
		logger.Fatal("please set your agent_reasoning_compute.base_url value in your configuration")
	}

	settings.BaseUrl = strings.TrimRight(settings.BaseUrl, "/")
	return &settings
}
