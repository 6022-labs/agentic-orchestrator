package models

import (
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Participant is anyone involved in a Conversation: a minted AI agent, a
// minted human, or a non-minted human (display-name only).
//
// Invariant (enforced at the app layer):
//   - Minted: CollectionAddress, CollectionTokenId, Role all non-nil.
//   - Non-minted: DisplayName non-nil; minted fields all nil.
//
// When minted, ID equals the agent-identity service's UUID. For non-minted
// participants, the broker generates the UUID.
type Participant struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`

	CollectionAddress *string                  `gorm:"type:text;uniqueIndex:idx_participants_collection" json:"collectionAddress,omitempty"`
	CollectionTokenId *string                  `gorm:"type:numeric(78,0);uniqueIndex:idx_participants_collection" json:"collectionTokenId,omitempty"`
	EnsDomain         *string                  `gorm:"type:text" json:"ensDomain,omitempty"`
	Role              *enums.ParticipantRole   `gorm:"type:text" json:"role,omitempty"`

	DisplayName *string `gorm:"type:text" json:"displayName,omitempty"`

	IsProcessing    bool       `gorm:"type:boolean;not null;default:false;index:idx_participants_processing,priority:1" json:"isProcessing"`
	LastProcessedAt *time.Time `gorm:"type:timestamptz;index:idx_participants_processing,priority:2" json:"lastProcessedAt,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updatedAt"`
}

func (entity *Participant) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}
	if entity.CreatedAt.IsZero() {
		entity.CreatedAt = time.Now()
	}
	return
}

func (entity *Participant) BeforeSave(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now()
	return
}
