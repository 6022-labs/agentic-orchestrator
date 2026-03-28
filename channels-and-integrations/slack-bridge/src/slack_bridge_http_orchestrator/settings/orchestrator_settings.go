package settings

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

const OrchestratorSettingsKey = "orchestrator"

type OrchestratorSettings struct {
	ApiBaseUrl            string `koanf:"api_base_url"`
	AgentComputeApiBaseUrl string `koanf:"agent_compute_api_base_url"`
}

func NewOrchestratorSettings(logger *zap.Logger, k *koanf.Koanf) *OrchestratorSettings {
	settings := OrchestratorSettings{}
	if err := k.Unmarshal(OrchestratorSettingsKey, &settings); err != nil {
		logger.Fatal("failed to unmarshal orchestrator settings", zap.Error(err))
	}

	settings.ApiBaseUrl = strings.TrimSpace(settings.ApiBaseUrl)
	if len(settings.ApiBaseUrl) == 0 {
		logger.Fatal("please set your orchestrator.api_base_url value in your configuration")
	}
	settings.ApiBaseUrl = strings.TrimSuffix(settings.ApiBaseUrl, "/")

	settings.AgentComputeApiBaseUrl = strings.TrimSpace(settings.AgentComputeApiBaseUrl)
	if len(settings.AgentComputeApiBaseUrl) == 0 {
		logger.Fatal("please set your orchestrator.agent_compute_api_base_url value in your configuration")
	}
	settings.AgentComputeApiBaseUrl = strings.TrimSuffix(settings.AgentComputeApiBaseUrl, "/")

	return &settings
}
