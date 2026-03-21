package models

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
)

type ExternalGetLLMConnectionSettingsResponse struct {
	Model    string                      `json:"model"`
	BaseUrl  *string                     `json:"baseUrl,omitempty"`
	Platform enums.LLMConnectionPlatform `json:"platform"`
}

func (response *ExternalGetLLMConnectionSettingsResponse) ToLLMConnectionSettings() *interfaces.LLMConnectionSettings {
	return &interfaces.LLMConnectionSettings{
		Model:    response.Model,
		BaseUrl:  response.BaseUrl,
		Platform: response.Platform,
	}
}
