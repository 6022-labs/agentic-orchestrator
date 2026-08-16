package models

import "github.com/google/uuid"

type SlackUser struct {
	Id       string `gorm:"primaryKey;type:string;not null" json:"id"`
	Name     string `gorm:"type:string;not null"            json:"name"`
	RealName string `gorm:"type:string;not null"            json:"realName"`

	WalletAddress *string `gorm:"type:string;unique" json:"walletAddress"`

	PairedHumanAgentCollectionAddress *string     `gorm:"type:string;uniqueIndex:idx_paired_slack_user_to_human_agent_collection"                                                                                                     json:"pairedHumanAgentCollectionAddress"`
	PairedHumanAgentCollectionTokenId *uint64     `gorm:"type:uint;uniqueIndex:idx_paired_slack_user_to_human_agent_collection"                                                                                                       json:"pairedHumanAgentCollectionTokenId"`
	PairedHumanAgent                  *SlackAgent `gorm:"foreignKey:PairedHumanAgentCollectionAddress,PairedHumanAgentCollectionTokenId;references:CollectionAddress,CollectionTokenId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"pairedHumanAgent"`

	// BrokerParticipantId is assigned by the conversation broker via
	// RegisterUser. SlackUserSynchronizer populates it eagerly the first time
	// it sees a Slack user (channel-member fan-out or message-author lookup)
	// and caches it here; subsequent IngestMessage calls reuse the cached
	// UUID without re-registering.
	BrokerParticipantId *uuid.UUID `gorm:"type:uuid" json:"brokerParticipantId,omitempty"`
}
