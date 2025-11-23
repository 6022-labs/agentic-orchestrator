package requests

import (
	"strconv"
	"strings"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type Message struct {
	Role           enums.AgentContextMessageRole     `json:"role"`
	Author         string                            `json:"author"`
	Content        string                            `json:"content"`
	Platform       enums.AgentContextMessagePlatform `json:"platform"`
	ConversationId *string                           `json:"conversationId" optional:"true"` // Not the same goal as the conversation id below.
}

type AddMessagesToAgentContextRequest struct {
	AgentId  uuid.UUID `json:"agentId"`
	Messages []Message `json:"messages"`

	// You can add messages to a context from a context key or a conversation id
	ContextKey     *string `json:"contextKey" optional:"true"`
	ConversationId *string `json:"conversationId" optional:"true"`
}

func (r *AddMessagesToAgentContextRequest) Validate() error {
	if r.AgentId == uuid.Nil {
		return errors.NewValidationError("agentId", "cannot be empty")
	}

	if len(r.Messages) == 0 {
		return errors.NewValidationError("messages", "cannot be empty")
	}

	for index, message := range r.Messages {
		if len(strings.TrimSpace(message.Content)) == 0 {
			return errors.NewValidationError("messages.content", "cannot be empty, check index "+strconv.Itoa(index))
		}

		if !enums.GetAllAgentContextMessageRoles().Contains(message.Role) {
			return errors.NewValidationError("messages.role", "must be one of the valid types"+enums.GetAllAgentContextMessageRoles().String()+", check index "+strconv.Itoa(index))
		}

		if !enums.GetAllAgentContextMessagePlatforms().Contains(message.Platform) {
			return errors.NewValidationError("messages.platform", "must be one of the valid types"+enums.GetAllAgentContextMessagePlatforms().String()+", check index "+strconv.Itoa(index))
		}
	}

	if r.ConversationId == nil && r.ContextKey == nil {
		return errors.NewValidationError("conversationId", "must be set if contextKey is not set")
	}

	if r.ConversationId != nil && r.ContextKey != nil {
		return errors.NewValidationError("conversationId", "must not be set if contextKey is set")
	}

	if r.ConversationId != nil {
		if len(strings.TrimSpace(*r.ConversationId)) == 0 {
			return errors.NewValidationError("conversationId", "cannot be empty")
		}
	}

	if r.ContextKey != nil {
		if len(strings.TrimSpace(*r.ContextKey)) == 0 {
			return errors.NewValidationError("contextKey", "cannot be empty")
		}
	}

	return nil
}

func (r *Message) ToAgentContextMessage(agentContextId uuid.UUID) *models.AgentContextMessage {
	return &models.AgentContextMessage{
		Role:           r.Role,
		Author:         r.Author,
		Content:        r.Content,
		Platform:       r.Platform,
		AgentContextId: agentContextId,
		ConversationId: r.ConversationId,
		CreatedAt:      time.Now(),
	}
}
