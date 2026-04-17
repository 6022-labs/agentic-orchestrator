package models

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
	"github.com/google/uuid"
)

type SlackAgent struct {
	CollectionAddress string          `gorm:"primaryKey;type:string;not null" json:"collectionAddress"`
	CollectionTokenId uint64          `gorm:"primaryKey;type:bigint;not null" json:"collectionTokenId"`
	EnsDomain         string          `gorm:"type:string;not null" json:"ensDomain"`
	Role              enums.AgentRole `gorm:"type:string;not null" json:"role"`
	IconImageCid      *string         `gorm:"type:string" json:"iconImageCid"`
	OwnerWallet       string          `gorm:"type:string;not null" json:"ownerWallet"`

	// BrokerParticipantId is the broker's participant UUID for this agent,
	// set by SlackAgentSynchronizer on first registration.
	BrokerParticipantId *uuid.UUID `gorm:"type:uuid" json:"brokerParticipantId,omitempty"`
}
