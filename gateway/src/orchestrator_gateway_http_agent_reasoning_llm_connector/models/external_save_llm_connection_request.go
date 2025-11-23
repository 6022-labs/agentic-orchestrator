package models

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/google/uuid"
)

type ExternalSaveLLMConnectionRequest struct {
	Model    string                      `json:"model"`
	ApiKey   *string                     `json:"apiKey"`
	AgentId  uuid.UUID                   `json:"agentId"`
	BaseUrl  *string                     `json:"baseUrl"`
	Platform enums.LLMConnectionPlatform `json:"platform"`

	AmazonPlatformRegion *string `json:"amazonPlatformRegion"`
}
