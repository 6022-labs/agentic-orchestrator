package models

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/compute/models/enums"
	"github.com/google/uuid"
)

type Agent struct {
	ID uuid.UUID `json:"id"`

	CollectionTokenId *uint64 `json:"collectionTokenId"`
	CollectionAddress string  `json:"collectionAddress"`
	ProposalId        *uint64 `json:"proposalId"`

	Name          string          `json:"name"`
	CloneOf       *uint64         `json:"cloneOf"`
	EnsDomain     string          `json:"ensDomain"`
	Role          enums.AgentRole `json:"role"`
	Creator       string          `json:"creator"`
	WalletAddress string          `json:"walletAddress"`
	Owner         *string         `json:"owner"`

	Images     map[string]string `json:"images"`
	Attributes map[string]string `json:"attributes"`

	IsInternal bool `json:"isInternal"`

	ImportationStatus *enums.AgentImportationStatus `json:"importationStatus"`
}
