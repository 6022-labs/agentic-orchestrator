package interfaces

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
	"github.com/google/uuid"
)

type LLMConnectionSettings struct {
	Model    string
	BaseUrl  *string
	Platform enums.LLMConnectionPlatform
}

type SaveLLMConnectionRequest struct {
	AgentId  uuid.UUID
	BaseUrl  *string
	ApiKey   *string
	Model    string
	Platform enums.LLMConnectionPlatform

	AmazonPlatformRegion *string
}

type LLMConnectionStatus struct {
	IsCompleted bool
	Status      *string
}

type LLMConnectionRequesterInterface interface {
	Remove(agentId uuid.UUID) error
	Save(request *SaveLLMConnectionRequest) error
	GetSettings(agentId uuid.UUID) (*LLMConnectionSettings, error)
	GetSetupStatus(agentId uuid.UUID) (*LLMConnectionStatus, error)
}
