package responses

import (
	"github.com/google/uuid"
)

type AgentResponse struct {
	ID uuid.UUID `json:"id"`

	Name       string            `json:"name"`
	Owner      *string           `json:"owner"`
	Creator    string            `json:"creator"`
	Images     map[string]string `json:"images"`
	CloneOf    *uint64           `json:"cloneOf"`
	Attributes map[string]string `json:"attributes"`
	EnsDomain  string            `json:"ensDomain"`
	Addresses  map[string]string `json:"addresses"`

	ProposalId        *uint64 `json:"proposalId"`
	CollectionTokenId *uint64 `json:"collectionTokenId"`
	CollectionAddress string  `json:"collectionAddress"`
}
