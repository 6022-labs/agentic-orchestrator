package repositories

import (
	"context"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
)

type AgentPromptTemplateRepositoryInterface interface {
	Save(ctx context.Context, entity *models.AgentPromptTemplate) (*models.AgentPromptTemplate, error)
	GetByAgentAndKey(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
		key enums.AgentPromptTemplateKey,
	) (*models.AgentPromptTemplate, error)
	ListByAgent(
		ctx context.Context,
		collectionAddress string,
		collectionTokenId uint64,
	) ([]*models.AgentPromptTemplate, error)
}
