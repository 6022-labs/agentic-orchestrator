package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	AgentPrivateKeyVaultPathTemplate = "agent/%s/private_key"
)

// Agent holds the schema definition for the Agent entity.
type Agent struct {
	ID uuid.UUID `gorm:"type:uuid;not null" json:"id"`

	CollectionTokenId *uint64 `gorm:"type:numeric" json:"collectionTokenId"`
	CollectionAddress string  `json:"collectionAddress"`
	ProposalId        *uint64 `gorm:"type:numeric" json:"proposalId"`

	Name          string          `gorm:"not null" json:"name"`
	CloneOf       *uint64         `gorm:"type:numeric" json:"cloneOf"`
	EnsDomain     string          `gorm:"not null" json:"ensDomain"`
	Role          enums.AgentRole `gorm:"type:text;not null" json:"role"`
	Creator       string          `gorm:"not null" json:"creator"`
	WalletAddress string          `gorm:"not null" json:"walletAddress"`
	Owner         *string         `json:"owner"`

	Images     map[string]string `gorm:"type:jsonb;serializer:json" json:"images"`
	Attributes map[string]string `gorm:"type:jsonb;serializer:json" json:"attributes"`


	UpdatedAt time.Time  `gorm:"type:timestamp;not null" json:"updatedAt"`
	CreatedAt time.Time  `gorm:"type:timestamp;not null" json:"createdAt"`
	DeletedAt *time.Time `gorm:"type:timestamp" json:"deletedAt"`
}

func (entity *Agent) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	entity.CreatedAt = time.Now()
	return
}

func (entity *Agent) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
