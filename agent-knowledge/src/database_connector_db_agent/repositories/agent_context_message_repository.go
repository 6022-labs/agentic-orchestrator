package repositories

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/db_clients"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentContextMessageRepository struct {
	*AbstractRepository[models.AgentContextMessage]
}

func NewAgentContextMessageRepository(client db_clients.AgentDatabaseClientInterface) *AgentContextMessageRepository {
	return &AgentContextMessageRepository{
		AbstractRepository: NewAbstractRepository[models.AgentContextMessage](client),
	}
}

func (r *AgentContextMessageRepository) SaveMultiple(entities []*models.AgentContextMessage) error {
	if len(entities) == 0 {
		return nil
	}
	err := r.client.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if err := tx.Save(entity).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *AgentContextMessageRepository) GetByAgentContextId(agentContextId uuid.UUID) ([]*models.AgentContextMessage, error) {
	var messages []*models.AgentContextMessage

	quotedCreatedAtColumn := quoteColumnName(r.client.GetDriver(), "created_at")
	quotedAgentContextIdColumn := quoteColumnName(r.client.GetDriver(), "agent_context_id")

	err := r.client.GetDB().Model(&models.AgentContextMessage{}).
		Where(quotedAgentContextIdColumn+" = ?", agentContextId).
		Order(fmt.Sprintf("agent_context_messages.%s %s", quotedCreatedAtColumn, "ASC")).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *AgentContextMessageRepository) GetByAgentContextKey(agentContextKey string) ([]*models.AgentContextMessage, error) {
	var messages []*models.AgentContextMessage

	quotedKeyColumn := quoteColumnName(r.client.GetDriver(), "key")
	quotedCreatedAtColumn := quoteColumnName(r.client.GetDriver(), "created_at")

	err := r.client.GetDB().Model(&models.AgentContextMessage{}).
		Joins("JOIN agent_contexts ON agent_contexts.id = agent_context_messages.agent_context_id").
		Where(fmt.Sprintf("agent_contexts.%s", quotedKeyColumn)+" = ?", agentContextKey).
		Order(fmt.Sprintf("agent_context_messages.%s ASC", quotedCreatedAtColumn)).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
