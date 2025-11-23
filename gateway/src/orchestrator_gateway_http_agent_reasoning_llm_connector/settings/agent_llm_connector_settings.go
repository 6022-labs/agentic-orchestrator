package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type AgentLLMConnectorSettings struct {
	BaseUrl string
}

func NewAgentLLMConnectorSettings(logger *zap.Logger) *AgentLLMConnectorSettings {
	agentLLMConnectorBaseUrl := os.Getenv("AGENT_LLM_CONNECTOR_BASE_URL")

	if len(strings.TrimSpace(agentLLMConnectorBaseUrl)) == 0 {
		logger.Fatal("AGENT_LLM_CONNECTOR_BASE_URL is not set")
	}

	agentLLMConnectorBaseUrl = strings.TrimRight(agentLLMConnectorBaseUrl, "/")

	return &AgentLLMConnectorSettings{
		BaseUrl: agentLLMConnectorBaseUrl,
	}
}
