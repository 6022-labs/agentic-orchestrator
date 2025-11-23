package settings

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type VerifyWalletOwnershipSettings struct {
	OrchestratorFrontBaseUrl string
}

func NewVerifyWalletOwnershipSettings(logger *zap.Logger) *VerifyWalletOwnershipSettings {
	orchestratorFrontBaseUrl := os.Getenv("ORCHESTRATOR_FRONT_BASE_URL")
	if len(strings.TrimSpace(orchestratorFrontBaseUrl)) == 0 {
		logger.Fatal("ORCHESTRATOR_FRONT_BASE_URL is not set")
	}

	orchestratorFrontBaseUrl = strings.TrimSuffix(orchestratorFrontBaseUrl, "/")

	return &VerifyWalletOwnershipSettings{
		OrchestratorFrontBaseUrl: orchestratorFrontBaseUrl,
	}
}
