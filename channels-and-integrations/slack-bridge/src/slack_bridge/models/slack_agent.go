package models

import (
	"time"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge/models/enums"
)

type SlackAgent struct {
	CollectionAddress string          `gorm:"primaryKey;type:string;not null" json:"collectionAddress"`
	CollectionTokenId uint64          `gorm:"primaryKey;type:bigint;not null" json:"collectionTokenId"`
	Name              string          `gorm:"type:string;not null" json:"name"`
	EnsDomain         string          `gorm:"type:string;not null" json:"ensDomain"`
	Role              enums.AgentRole `gorm:"type:string;not null" json:"role"`
	IconImageCid      *string         `gorm:"type:string" json:"iconImageCid"`
	OwnerWallet       string          `gorm:"type:string;not null" json:"ownerWallet"`
	IsProcessing      bool            `gorm:"type:boolean;not null" json:"isProcessing"`
	LastProcessedAt   *time.Time      `gorm:"type:timestamp" json:"lastProcessedAt"`
}
