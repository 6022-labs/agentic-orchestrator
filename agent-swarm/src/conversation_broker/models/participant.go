package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Participant is anyone involved in a Conversation: an agent (AI clone,
// expert, facilitator, or human paired to an NFT) or a user (display-name
// only).
//
// Invariant (enforced at the app layer):
//   - Agent: CollectionAddress and CollectionTokenId non-nil.
//   - User: DisplayName non-nil; on-chain identity fields all nil.
//
// For agents, ID equals the agent-identity service's UUID. For users, the
// broker generates the UUID.
type Participant struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`

	CollectionAddress *string `gorm:"type:text;uniqueIndex:idx_participants_collection"    json:"collectionAddress,omitempty"`
	CollectionTokenId *uint64 `gorm:"type:numeric;uniqueIndex:idx_participants_collection" json:"collectionTokenId,omitempty"`
	EnsDomain         *string `gorm:"type:text"                                            json:"ensDomain,omitempty"`

	DisplayName *string `gorm:"type:text" json:"displayName,omitempty"`

	IsProcessing    bool       `gorm:"type:boolean;not null;default:false;index:idx_participants_processing,priority:1" json:"isProcessing"`
	LastProcessedAt *time.Time `gorm:"type:timestamptz;index:idx_participants_processing,priority:2"                    json:"lastProcessedAt,omitempty"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updatedAt"`
}

// IsAgent reports whether the participant has an on-chain identity rather than being a user.
func (entity *Participant) IsAgent() bool {
	return entity.CollectionAddress != nil && entity.CollectionTokenId != nil
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
