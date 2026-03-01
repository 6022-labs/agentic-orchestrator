package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgentPromptTemplate struct {
	ID      uuid.UUID                    `gorm:"type:char(36);primaryKey"`
	Key     enums.AgentPromptTemplateKey `gorm:"type:varchar(64);not null;unique"`
	Content string                       `gorm:"type:text;not null"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (entity *AgentPromptTemplate) BeforeCreate(tx *gorm.DB) error {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}

	return nil
}

func (entity *AgentPromptTemplate) BeforeSave(tx *gorm.DB) error {
	entity.UpdatedAt = time.Now()

	return nil
}
