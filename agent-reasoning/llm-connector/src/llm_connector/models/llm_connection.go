package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/llm-connector/src/llm_connector/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	LLMSecretVaultPathTemplate = "agent/%s/llm_secret"
)

// LLMConnection holds the schema definition for the LLMConnection entity.
type LLMConnection struct {
	ID uuid.UUID `gorm:"type:uuid;not null" json:"id"`

	Model    string                      `gorm:"type:string" json:"model"`
	BaseUrl  *string                     `gorm:"type:string" json:"baseUrl"`
	AgentID  uuid.UUID                   `gorm:"type:uuid;not null" json:"agentId"`
	Platform enums.LLMConnectionPlatform `gorm:"type:string;not null" json:"platform"`

	SetupStatus enums.LLMConnectionSetupStatus `gorm:"type:string;not null" json:"setupStatus"`

	UpdatedAt time.Time  `gorm:"type:timestamp;not null" json:"updatedAt"`
	CreatedAt time.Time  `gorm:"type:timestamp;not null" json:"createdAt"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"deletedAt"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (entity *LLMConnection) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	entity.CreatedAt = time.Now()
	return
}

// BeforeSave will update the UpdatedAt field.
func (entity *LLMConnection) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
