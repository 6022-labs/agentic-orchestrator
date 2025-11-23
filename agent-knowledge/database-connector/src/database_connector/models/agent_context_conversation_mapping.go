package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentContextConversationMapping struct {
	ID             uuid.UUID `gorm:"type:char(36);primaryKey"`
	ConversationId string    `gorm:"size:1024;not null"`

	CreatedAt time.Time `gorm:"not null"`

	AgentContextId uuid.UUID    `gorm:"type:char(36);index;not null"`
	AgentContext   AgentContext `gorm:"foreignKey:AgentContextId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (entity *AgentContextConversationMapping) BeforeCreate(tx *gorm.DB) error {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}

	return nil
}
