package responses

import "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/llm_connector/models/enums"

type GetLLMConnectionSettingsResponse struct {
	BaseUrl  *string                     `json:"baseUrl,omitempty"`
	Model    string                      `json:"model"`
	Platform enums.LLMConnectionPlatform `json:"platform"`
}
