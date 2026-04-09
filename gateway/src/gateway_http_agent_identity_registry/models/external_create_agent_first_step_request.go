package models

import (
	"mime/multipart"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
)

type ExternalCreateAgentFirstStepRequest struct {
	AuthenticatedWallet string
	CollectionAddress   string
	Name                string
	Role                enums.AgentRole
	WalletAddress       *string
	Mission             string
	CloneOfAgentId      *uint64
	Image               *multipart.FileHeader
}
