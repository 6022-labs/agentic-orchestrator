package responses

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/google/uuid"
)

type AgentPromptTemplateResponse struct {
	ID                uuid.UUID                    `json:"id"`
	CollectionAddress string                       `json:"collectionAddress"`
	CollectionTokenId uint64                       `json:"collectionTokenId"`
	Key               enums.AgentPromptTemplateKey `json:"key"`
	Content           string                       `json:"content"`
	CreatedAt         time.Time                    `json:"createdAt"`
	UpdatedAt         time.Time                    `json:"updatedAt"`
}

func AgentPromptTemplateResponseFromModel(template *models.AgentPromptTemplate) *AgentPromptTemplateResponse {
	if template == nil {
		return nil
	}
	return &AgentPromptTemplateResponse{
		ID:                template.ID,
		CollectionAddress: template.CollectionAddress,
		CollectionTokenId: template.CollectionTokenId,
		Key:               template.Key,
		Content:           template.Content,
		CreatedAt:         template.CreatedAt,
		UpdatedAt:         template.UpdatedAt,
	}
}
