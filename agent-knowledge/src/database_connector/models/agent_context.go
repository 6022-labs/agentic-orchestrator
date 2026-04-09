package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentContext struct {
	ID  uuid.UUID `gorm:"type:char(36);primaryKey"`
	Key string    `gorm:"size:256;not null;uniqueIndex"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	Messages             []AgentContextMessage             `gorm:"foreignKey:AgentContextId"`
	ConversationMappings []AgentContextConversationMapping `gorm:"foreignKey:AgentContextId"`
}

func (entity *AgentContext) BeforeCreate(tx *gorm.DB) error {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}

	return nil
}

func (entity *AgentContext) BeforeSave(tx *gorm.DB) error {
	entity.UpdatedAt = time.Now()

	return nil
}
