package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type OrchestratorSettings struct {
	ApiBaseUrl string
}

func NewOrchestratorSettings(logger *zap.Logger) *OrchestratorSettings {
	apiBaseUrl := os.Getenv("ORCHESTRATOR_API_BASE_URL")
	if len(strings.TrimSpace(apiBaseUrl)) == 0 {
		logger.Fatal("ORCHESTRATOR_API_BASE_URL is not set")
	}

	apiBaseUrl = strings.TrimSuffix(apiBaseUrl, "/")

	return &OrchestratorSettings{
		ApiBaseUrl: apiBaseUrl,
	}
}
