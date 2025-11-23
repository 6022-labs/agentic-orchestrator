package models

import (
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/models/enums"
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
	CollectionAddress string  `gorm:"type:string" json:"collectionAddress"`
	ProposalId        *uint64 `gorm:"type:numeric" json:"proposalId"`

	Name          string          `gorm:"type:string;not null" json:"name"`
	CloneOf       *uint64         `gorm:"type:numeric" json:"cloneOf"`
	EnsDomain     string          `gorm:"type:string;not null" json:"ensDomain"`
	Role          enums.AgentRole `gorm:"type:string;not null" json:"role"`
	Creator       string          `gorm:"type:string;not null" json:"creator"`
	WalletAddress string          `gorm:"type:string;not null" json:"walletAddress"`
	Owner         *string         `gorm:"type:string" json:"owner"`

	Images     map[string]string `gorm:"type:jsonb;serializer:json" json:"images"`
	Attributes map[string]string `gorm:"type:jsonb;serializer:json" json:"attributes"`

	IsInternal bool `gorm:"type:boolean;not null" json:"isInternal"`

	ImportationStatus *enums.AgentImportationStatus `gorm:"type:string" json:"importationStatus"`

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
