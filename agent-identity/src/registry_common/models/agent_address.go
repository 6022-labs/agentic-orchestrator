package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgentAddress is one of the agent's typed on-chain addresses, one row per chain family keyed by (agent_id, type).
type AgentAddress struct {
	AgentId uuid.UUID         `gorm:"type:uuid;primaryKey" json:"agentId"`
	Type    enums.AddressType `gorm:"type:text;primaryKey" json:"type"`
	Address string            `gorm:"not null" json:"address"`

	UpdatedAt time.Time `gorm:"type:timestamp;not null" json:"updatedAt"`
	CreatedAt time.Time `gorm:"type:timestamp;not null" json:"createdAt"`
}

func (AgentAddress) TableName() string {
	return "agent_addresses"
}

func (entity *AgentAddress) BeforeCreate(tx *gorm.DB) (err error) {
	entity.CreatedAt = time.Now()
	return
}

func (entity *AgentAddress) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
