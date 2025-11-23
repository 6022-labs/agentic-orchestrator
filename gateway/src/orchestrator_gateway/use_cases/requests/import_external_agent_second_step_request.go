package requests

import (
	"strings"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
)

type ImportExternalAgentSecondStepRequest struct {
	AgentByCollectionTokenIdRequest
	LLMBaseUrl     *string                        `json:"llmBaseUrl"`
	LLMApiKey      *string                        `json:"llmApiKey"`
	LLMModel       string                         `json:"llmModel"`
	LLMPlatform    enums.LLMConnectionPlatform    `json:"llmPlatform"`
	DatabaseHost   string                         `json:"databaseHost"`
	DatabasePort   uint16                         `json:"databasePort"`
	DatabaseDriver enums.DatabaseConnectionDriver `json:"databaseDriver"`
}

func (r *ImportExternalAgentSecondStepRequest) Validate() error {
	if len(strings.TrimSpace(r.LLMModel)) <= 2 {
		return errors.NewValidationError("llmModel", "must be at least 2 characters long")
	}

	if !enums.GetAllLLMConnectionPlatforms().Contains(r.LLMPlatform) {
		return errors.NewValidationError("llmPlatform", "must be one of the valid platforms: "+enums.GetAllLLMConnectionPlatforms().String())
	}

	if !enums.GetAllDatabaseConnectionDrivers().Contains(r.DatabaseDriver) {
		return errors.NewValidationError("databaseDriver", "must be one of the valid database driver: "+enums.GetAllDatabaseConnectionDrivers().String())
	}

	if enums.GetAllLocalLLMConnectionPlatforms().Contains(r.LLMPlatform) {
		if r.LLMBaseUrl == nil || len(strings.TrimSpace(*r.LLMBaseUrl)) == 0 {
			return errors.NewValidationError("llmBaseUrl", "must be provided for the specified platform: "+r.LLMPlatform.String())
		}
	}

	if enums.GetAllCloudLLMConnectionPlatforms().Contains(r.LLMPlatform) {
		if r.LLMApiKey == nil || len(strings.TrimSpace(*r.LLMApiKey)) == 0 {
			return errors.NewValidationError("llmApiKey", "must be provided for the specified platform: "+r.LLMPlatform.String())
		}
	}

	return r.EnsureAuthenticatedWallet()
}
