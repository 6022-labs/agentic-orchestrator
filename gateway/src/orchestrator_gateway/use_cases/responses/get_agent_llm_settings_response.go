package responses

import "github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"

type GetAgentLLMSettingsResponse struct {
	Model    string                      `json:"model"`
	BaseUrl  *string                     `json:"baseUrl,omitempty"`
	Platform enums.LLMConnectionPlatform `json:"platform"`
}
