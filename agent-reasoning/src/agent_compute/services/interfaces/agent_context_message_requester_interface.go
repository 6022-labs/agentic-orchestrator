package interfaces

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/agent_compute/models/enums"
	"github.com/google/uuid"
)

type AddMessageToAgentContextRequestMessage struct {
	Content        string
	Author         string
	Role           enums.AgentContextMessageRole
	Platform       enums.AgentContextMessagePlatform
	ConversationId *string
}

type AddMessageToAgentContextRequest struct {
	AgentId        uuid.UUID
	ConversationId *string
	ContextKey     *string
	Messages       []*AddMessageToAgentContextRequestMessage
}

type AgentContextMessageRequesterInterface interface {
	AddMessagesToAgentContext(request *AddMessageToAgentContextRequest) error
	GetAgentContextMessages(agentId uuid.UUID, key string) ([]*models.AgentContextMessage, error)
}
