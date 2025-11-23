package models

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-knowledge/database-connector/src/database_connector/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentContextMessage struct {
	ID             uuid.UUID                         `gorm:"type:varchar(36);primaryKey"`
	Content        string                            `gorm:"type:text;not null"`
	Author         string                            `gorm:"type:varchar(64);not null"`
	Platform       enums.AgentContextMessagePlatform `gorm:"type:varchar(64);not null"`
	Role           enums.AgentContextMessageRole     `gorm:"type:varchar(64);not null"`
	Hash           string                            `gorm:"type:varchar(128);not null"`
	ConversationId *string                           `gorm:"type:varchar(128);optional"`

	CreatedAt time.Time `gorm:"not null"`

	AgentContextId uuid.UUID    `gorm:"type:char(36);index;not null"`
	AgentContext   AgentContext `gorm:"foreignKey:AgentContextId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (entity *AgentContextMessage) BeforeCreate(tx *gorm.DB) error {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}

	return nil
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
