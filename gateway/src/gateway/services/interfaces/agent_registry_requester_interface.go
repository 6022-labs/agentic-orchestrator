package interfaces

import (
	"mime/multipart"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type ListRegistryAgentsRequest struct {
	Page            *uint
	ElementsPerPage *uint
	SortBy          *string
	SortDirection   *string
	IsInternal      *bool
}

type CreateRegistryAgentFirstStepRequest struct {
	AuthenticatedWallet common.Address
	CollectionAddress   string
	Name                string
	Role                enums.AgentRole
	WalletAddress       *string
	Mission             string
	CloneOfAgentId      *uint64
	Image               *multipart.FileHeader
}

type CreateRegistryAgentSecondStepRequest struct {
	AuthenticatedWallet common.Address
	AgentId             uuid.UUID
}

type AgentRegistryRequesterInterface interface {
	ListAgents(request *ListRegistryAgentsRequest) (*responses.ListAgentsResponse, error)
	GetAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*responses.AgentResponse, error)
	ListAgentCollections() ([]*responses.AgentCollectionResponse, error)
	IsInternalAgent(agentCollectionAddress string, agentCollectionTokenId uint64) (*responses.IsInternalAgentResponse, error)
	CreateAgentFirstStep(request *CreateRegistryAgentFirstStepRequest) (*responses.AgentResponse, error)
	CreateAgentSecondStep(request *CreateRegistryAgentSecondStepRequest) (*responses.CreateAgentSecondStepResponse, error)
	GetMintedAgentByEnsDomain(ensDomain string) (*responses.AgentResponse, error)
	GetAgentsByEnsDomains(ensDomains []string) ([]*responses.AgentResponse, error)
}
