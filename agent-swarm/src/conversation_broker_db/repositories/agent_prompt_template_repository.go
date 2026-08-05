package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
)

type AgentPromptTemplateRepository struct {
	AbstractRepository[models.AgentPromptTemplate]
}

func NewAgentPromptTemplateRepository(client db_clients.DatabaseClientInterface) *AgentPromptTemplateRepository {
	return &AgentPromptTemplateRepository{
		AbstractRepository: *NewAbstractRepository[models.AgentPromptTemplate](client),
	}
}

func (r *AgentPromptTemplateRepository) GetByAgentAndKey(
	collectionAddress string,
	collectionTokenId uint64,
	key enums.AgentPromptTemplateKey,
) (*models.AgentPromptTemplate, error) {
	var entity models.AgentPromptTemplate
	err := r.client.GetDB().
		Where("collection_address = ?", collectionAddress).
		Where("collection_token_id = ?", collectionTokenId).
		Where("key = ?", key).
		First(&entity).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *AgentPromptTemplateRepository) ListByAgent(
	collectionAddress string,
	collectionTokenId uint64,
) ([]*models.AgentPromptTemplate, error) {
	var entities []*models.AgentPromptTemplate
	err := r.client.GetDB().
		Where("collection_address = ?", collectionAddress).
		Where("collection_token_id = ?", collectionTokenId).
		Order("key ASC").
		Find(&entities).
		Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
