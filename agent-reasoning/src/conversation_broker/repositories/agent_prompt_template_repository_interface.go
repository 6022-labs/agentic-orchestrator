package repositories

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
)

type AgentPromptTemplateRepositoryInterface interface {
	Save(entity *models.AgentPromptTemplate) (*models.AgentPromptTemplate, error)
	GetByAgentAndKey(collectionAddress string, collectionTokenId uint64, key enums.AgentPromptTemplateKey) (*models.AgentPromptTemplate, error)
	ListByAgent(collectionAddress string, collectionTokenId uint64) ([]*models.AgentPromptTemplate, error)
}
