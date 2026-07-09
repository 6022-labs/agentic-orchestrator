package models

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/google/uuid"
)

type ExternalAgentResponse struct {
	ID uuid.UUID `json:"id"`

	Name          string            `json:"name"`
	Role          enums.AgentRole   `json:"role"`
	Owner         *string           `json:"owner"`
	Creator       string            `json:"creator"`
	Images        map[string]string `json:"images"`
	CloneOf       *uint64           `json:"cloneOf"`
	Attributes    map[string]string `json:"attributes"`
	EnsDomain     string            `json:"ensDomain"`
	WalletAddress string            `json:"walletAddress"`

	ProposalId        *uint64 `json:"proposalId"`
	CollectionTokenId *uint64 `json:"collectionTokenId"`
	CollectionAddress string  `json:"collectionAddress"`
}
