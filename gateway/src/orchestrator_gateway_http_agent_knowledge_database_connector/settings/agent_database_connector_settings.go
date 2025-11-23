package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type AgentDatabaseConnectorSettings struct {
	BaseUrl string
}

func NewAgentDatabaseConnectorSettings(logger *zap.Logger) *AgentDatabaseConnectorSettings {
	agentDatabaseConnectorBaseUrl := os.Getenv("AGENT_DATABASE_CONNECTOR_BASE_URL")

	if len(strings.TrimSpace(agentDatabaseConnectorBaseUrl)) == 0 {
		logger.Fatal("AGENT_DATABASE_CONNECTOR_BASE_URL is not set")
	}

	// Trim any trailing slashes from the BaseUrl
	agentDatabaseConnectorBaseUrl = strings.TrimRight(agentDatabaseConnectorBaseUrl, "/")

	return &AgentDatabaseConnectorSettings{
		BaseUrl: agentDatabaseConnectorBaseUrl,
	}
}
