package interfaces

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models"
	"github.com/google/uuid"
)

type AgentContextConversationMappingRequesterInterface interface {
	GetLinkedAgentContextKeyOfConversation(agentId uuid.UUID, conversationId string) (*models.AgentContext, *bool, error)
	LinkConversationToAgentContext(agentId uuid.UUID, conversationId string, key string) error
}
