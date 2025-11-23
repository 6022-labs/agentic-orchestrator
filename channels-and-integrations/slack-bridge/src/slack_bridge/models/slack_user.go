package models

type SlackUser struct {
	Id       string `gorm:"primaryKey;type:string;not null" json:"id"`
	Name     string `gorm:"type:string;not null" json:"name"`
	RealName string `gorm:"type:string;not null" json:"realName"`

	WalletAddress *string `gorm:"type:string;unique" json:"walletAddress"`

	PairedHumanAgentCollectionAddress *string     `gorm:"type:string;uniqueIndex:idx_paired_slack_user_to_human_agent_collection" json:"pairedHumanAgentCollectionAddress"`
	PairedHumanAgentCollectionTokenId *uint64     `gorm:"type:uint;uniqueIndex:idx_paired_slack_user_to_human_agent_collection" json:"pairedHumanAgentCollectionTokenId"`
	PairedHumanAgent                  *SlackAgent `gorm:"foreignKey:PairedHumanAgentCollectionAddress,PairedHumanAgentCollectionTokenId;references:CollectionAddress,CollectionTokenId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"pairedHumanAgent"`
}
