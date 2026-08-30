package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConversationMemoryBinding pins an agent participant to the node-side session its turns run on.
type ConversationMemoryBinding struct {
	ConversationID uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"conversationId"`
	ParticipantID  uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"participantId"`

	MemorySessionID uuid.UUID `gorm:"type:uuid;not null"        json:"memorySessionId"`
	BoundAt         time.Time `gorm:"type:timestamptz;not null" json:"boundAt"`

	Conversation *Conversation `gorm:"foreignKey:ConversationID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Participant  *Participant  `gorm:"foreignKey:ParticipantID;references:ID;constraint:OnDelete:CASCADE"  json:"-"`
}

func (entity *ConversationMemoryBinding) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.BoundAt.IsZero() {
		entity.BoundAt = time.Now()
	}
	return
}
