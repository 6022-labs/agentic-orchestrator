package responses

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
	"github.com/google/uuid"
)

type GetAgentContextMessageResponse struct {
	Hash           string                            `json:"hash"`
	Role           enums.AgentContextMessageRole     `json:"role"`
	Author         string                            `json:"author"`
	Content        string                            `json:"content"`
	Platform       enums.AgentContextMessagePlatform `json:"platform"`
	CreatedAt      time.Time                         `json:"createdAt"`
	ConversationId *string                           `json:"conversationId"`
	AgentContextId uuid.UUID                         `json:"agentContextId"`
}
