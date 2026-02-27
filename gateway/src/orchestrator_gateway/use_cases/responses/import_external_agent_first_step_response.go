package responses

import "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"

type ImportExternalAgentFirstStepResponse struct {
	Name        string                      `json:"name"`
	Role        string                      `json:"role"`
	LLMModel    string                      `json:"llmModel"`
	LLMPlatform enums.LLMConnectionPlatform `json:"llmPlatform"`
}
