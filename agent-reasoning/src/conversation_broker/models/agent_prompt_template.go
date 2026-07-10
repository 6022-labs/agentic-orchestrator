package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgentPromptTemplate is an agent-scoped prompt override; when no row exists
// for a key the processor falls back to the built-in role default.
type AgentPromptTemplate struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`

	CollectionAddress string                       `gorm:"type:text;not null;uniqueIndex:idx_agent_prompt_templates_agent_key,priority:1" json:"collectionAddress"`
	CollectionTokenId uint64                       `gorm:"type:numeric;not null;uniqueIndex:idx_agent_prompt_templates_agent_key,priority:2" json:"collectionTokenId"`
	Key               enums.AgentPromptTemplateKey `gorm:"type:text;not null;uniqueIndex:idx_agent_prompt_templates_agent_key,priority:3" json:"key"`
	Content           string                       `gorm:"type:text;not null" json:"content"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updatedAt"`
}

func (entity *AgentPromptTemplate) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}
	return
}

func (entity *AgentPromptTemplate) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
