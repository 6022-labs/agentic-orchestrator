package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/orchestrator/models/enums"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ExternalGetAgentContextMessageResponse struct {
	Hash           string                            `json:"hash"`
	Role           enums.AgentContextMessageRole     `json:"role"`
	Author         string                            `json:"author"`
	Content        string                            `json:"content"`
	Platform       enums.AgentContextMessagePlatform `json:"platform"`
	CreatedAt      time.Time                         `json:"createdAt"`
	ConversationId *string                           `json:"conversationId"`
	AgentContextId uuid.UUID                         `json:"agentContextId"`
}

type ExternalGetAgentContextMessageResponses []*ExternalGetAgentContextMessageResponse

func (e *ExternalGetAgentContextMessageResponses) ToAgentContextMessages() []*models.AgentContextMessage {
	return lo.Map(*e, func(message *ExternalGetAgentContextMessageResponse, _ int) *models.AgentContextMessage {
		if message == nil {
			return nil
		}

		return &models.AgentContextMessage{
			Hash:           message.Hash,
			Role:           message.Role,
			Author:         message.Author,
			Content:        message.Content,
			Platform:       message.Platform,
			CreatedAt:      message.CreatedAt,
			ConversationId: message.ConversationId,
			AgentContextId: message.AgentContextId,
		}
	})
}
