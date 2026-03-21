package models

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/google/uuid"
)

// Replicate of the model in database connector without id (not needed here)
type AgentContextMessage struct {
	Content        string
	Author         string
	Platform       enums.AgentContextMessagePlatform
	Role           enums.AgentContextMessageRole
	Hash           string
	ConversationId *string

	CreatedAt time.Time

	AgentContextId uuid.UUID
}

func (entity *AgentContextMessage) EnsureHash(agentContextId uuid.UUID) error {
	if entity.Hash != "" {
		return nil
	}

	entity.AgentContextId = agentContextId

	return entity.ComputeHash()
}

func (entity *AgentContextMessage) ComputeHash() error {
	b, err := json.Marshal(struct {
		Role           enums.AgentContextMessageRole     `json:"role"`
		Author         string                            `json:"author"`
		Content        string                            `json:"content"`
		Platform       enums.AgentContextMessagePlatform `json:"platform"`
		ConversationId *string                           `json:"conversationId,omitempty"`
		AgentContextId uuid.UUID                         `json:"agentContextId"`
	}{
		Role:           entity.Role,
		Author:         entity.Author,
		Content:        entity.Content,
		Platform:       entity.Platform,
		ConversationId: entity.ConversationId,
		AgentContextId: entity.AgentContextId,
	})
	if err != nil {
		return err
	}

	// SHA-512 hash
	hash := sha512.Sum512(b)

	entity.Hash = hex.EncodeToString(hash[:])

	return nil
}
